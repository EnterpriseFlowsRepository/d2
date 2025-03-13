package shape

import (
	"math"

	"oss.terrastruct.com/d2/lib/geo"
	"oss.terrastruct.com/d2/lib/svg"
	"oss.terrastruct.com/util-go/go2"
)

// Optimal value for circular arc approximation with cubic bezier curves
const kCircleApprox = 0.5522847498307936 // 4*(math.Sqrt(2)-1)/3

type shapeC4Person struct {
	*baseShape
}

func NewC4Person(box *geo.Box) Shape {
	shape := shapeC4Person{
		baseShape: &baseShape{
			Type: C4_PERSON_TYPE,
			Box:  box,
		},
	}
	shape.FullShape = go2.Pointer(Shape(shape))
	return shape
}

const (
	C4_PERSON_AR_LIMIT = 1.5
)

func (s shapeC4Person) GetInnerBox() *geo.Box {
	width := s.Box.Width
	height := s.Box.Height

	headRadius := width * 0.22
	headCenterY := height * 0.18
	bodyTop := headCenterY + headRadius*0.8

	// Use a small fixed percentage instead of the full corner radius
	horizontalPadding := width * 0.05 // 5% padding

	tl := s.Box.TopLeft.Copy()
	tl.X += horizontalPadding

	// Add vertical padding
	tl.Y += bodyTop + height*0.03

	// Width minus padding on both sides
	innerWidth := width - (horizontalPadding * 2)

	// Add bottom padding
	innerHeight := height - (tl.Y - s.Box.TopLeft.Y) - (height * 0.03)

	return geo.NewBox(tl, innerWidth, innerHeight)
}

func bodyPath(box *geo.Box) *svg.SvgPathContext {
	width := box.Width
	height := box.Height

	pc := svg.NewSVGPathContext(box.TopLeft, 1, 1)

	headRadius := width * 0.22
	headCenterY := height * 0.18
	bodyTop := headCenterY + headRadius*0.8
	bodyWidth := width
	bodyHeight := height - bodyTop
	bodyLeft := 0
	// Ensure cornerRadius is constrained to a portion of the shortest dimension
	// This prevents distorted corners when width is large compared to height
	cornerRadius := math.Min(width*0.175, bodyHeight*0.25)

	pc.StartAt(pc.Absolute(float64(bodyLeft), bodyTop+cornerRadius))

	pc.C(true, 0, -kCircleApprox*cornerRadius, kCircleApprox*cornerRadius, -cornerRadius, cornerRadius, -cornerRadius)
	pc.H(true, bodyWidth-2*cornerRadius)
	pc.C(true, kCircleApprox*cornerRadius, 0, cornerRadius, kCircleApprox*cornerRadius, cornerRadius, cornerRadius)
	pc.V(true, bodyHeight-2*cornerRadius)
	pc.C(true, 0, kCircleApprox*cornerRadius, -kCircleApprox*cornerRadius, cornerRadius, -cornerRadius, cornerRadius)
	pc.H(true, -(bodyWidth - 2*cornerRadius))
	pc.C(true, -kCircleApprox*cornerRadius, 0, -cornerRadius, -kCircleApprox*cornerRadius, -cornerRadius, -cornerRadius)
	pc.Z()

	return pc
}

func headPath(box *geo.Box) *svg.SvgPathContext {
	width := box.Width
	height := box.Height

	pc := svg.NewSVGPathContext(box.TopLeft, 1, 1)

	headRadius := width * 0.22
	headCenterX := width / 2
	headCenterY := height * 0.18

	pc.StartAt(pc.Absolute(headCenterX, headCenterY-headRadius))

	pc.C(false,
		headCenterX+headRadius*kCircleApprox, headCenterY-headRadius,
		headCenterX+headRadius, headCenterY-headRadius*kCircleApprox,
		headCenterX+headRadius, headCenterY)

	pc.C(false,
		headCenterX+headRadius, headCenterY+headRadius*kCircleApprox,
		headCenterX+headRadius*kCircleApprox, headCenterY+headRadius,
		headCenterX, headCenterY+headRadius)

	pc.C(false,
		headCenterX-headRadius*kCircleApprox, headCenterY+headRadius,
		headCenterX-headRadius, headCenterY+headRadius*kCircleApprox,
		headCenterX-headRadius, headCenterY)

	pc.C(false,
		headCenterX-headRadius, headCenterY-headRadius*kCircleApprox,
		headCenterX-headRadius*kCircleApprox, headCenterY-headRadius,
		headCenterX, headCenterY-headRadius)

	return pc
}

func (s shapeC4Person) Perimeter() []geo.Intersectable {
	width := s.Box.Width
	height := s.Box.Height

	bodyPerimeter := bodyPath(s.Box).Path

	headRadius := width * 0.22
	headCenterX := s.Box.TopLeft.X + width/2
	headCenterY := s.Box.TopLeft.Y + height*0.18
	headCenter := geo.NewPoint(headCenterX, headCenterY)

	headEllipse := geo.NewEllipse(headCenter, headRadius, headRadius)

	return append(bodyPerimeter, headEllipse)
}

func (s shapeC4Person) GetSVGPathData() []string {
	return []string{
		bodyPath(s.Box).PathData(),
		headPath(s.Box).PathData(),
	}
}

func (s shapeC4Person) GetDimensionsToFit(width, height, paddingX, paddingY float64) (float64, float64) {
	contentWidth := width + paddingX
	contentHeight := height + paddingY

	// Account for 10% total horizontal padding (5% on each side)
	totalWidth := contentWidth / 0.9
	headRadius := totalWidth * 0.22
	headCenterY := totalWidth * 0.18
	bodyTop := headCenterY + headRadius*0.8

	// Include vertical padding from GetInnerBox
	verticalPadding := totalWidth * 0.06 // 3% top + 3% bottom
	totalHeight := contentHeight + bodyTop + verticalPadding

	minHeight := totalWidth * 1.2
	if totalHeight < minHeight {
		totalHeight = minHeight
	}

	totalWidth, totalHeight = LimitAR(totalWidth, totalHeight, C4_PERSON_AR_LIMIT)
	return math.Ceil(totalWidth), math.Ceil(totalHeight)
}

func (s shapeC4Person) GetDefaultPadding() (paddingX, paddingY float64) {
	return 20, defaultPadding * 1.5
}
