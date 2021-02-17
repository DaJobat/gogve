package sweepline

import (
	"fmt"
	"math"

	"github.com/DaJobat/gogve/tree"
	"github.com/DaJobat/gogve/util"
)

type SweeplineState interface {
	CurrentSweepline() float64
}

type LineSegment struct {
	startPoint util.FVec
	endPoint   util.FVec
	gradient   float64
	yIntercept float64 //obviously this is not a "real" y intercept,
	// but where the line would intercept if it continued infinitely (as a line)
}

func NewLineSegment(startPoint, endPoint util.FVec) *LineSegment {
	if startPoint.Degree() != 2 || endPoint.Degree() != 2 {
		panic("not a 2 vec for line")
	}

	var ls LineSegment
	if startPoint.X() < endPoint.X() {
		ls = LineSegment{
			startPoint: startPoint,
			endPoint:   endPoint,
		}
	} else {
		ls = LineSegment{
			startPoint: endPoint,
			endPoint:   startPoint,
		}
	}
	var cx, cy float64
	if endPoint.X() == startPoint.X() {
		cx = 0
	} else {
		cx = endPoint.X() - startPoint.X()
	}

	if startPoint.Y() == endPoint.Y() {
		cy = 0
	} else {
		cy = endPoint.Y() - startPoint.Y()
	}
	if cx != 0 {
		ls.gradient = cy / cx
		if ls.gradient == 0 {
			ls.yIntercept = startPoint.Y()
		} else {
			ls.yIntercept = startPoint.Y() - (ls.gradient * startPoint.X())
		}
	} else {
		//Edge case for vertical lines
		ls.gradient = math.Inf(1)
		ls.yIntercept = math.Inf(1)
	}

	return &ls
}

func (l *LineSegment) Start() util.FVec {
	return l.startPoint
}

func (l *LineSegment) End() util.FVec {
	return l.endPoint
}

func (l *LineSegment) Gradient() float64 {
	return l.gradient
}

func (l *LineSegment) YIntercept() float64 {
	return l.yIntercept
}

//TODO: Find the intersection Point
func (l *LineSegment) Intersects(l1 *LineSegment) (intersects bool) {
	d1 := direction(l1.startPoint, l1.endPoint, l.startPoint)
	d2 := direction(l1.startPoint, l1.endPoint, l.endPoint)
	d3 := direction(l.startPoint, l.endPoint, l1.startPoint)
	d4 := direction(l.startPoint, l.endPoint, l1.endPoint)
	if ((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)) {
		return true
	} else if d1 == 0 && onSegment(l1.startPoint, l1.endPoint, l.startPoint) {
		return true
	} else if d2 == 0 && onSegment(l1.startPoint, l1.endPoint, l.endPoint) {
		return true
	} else if d3 == 0 && onSegment(l.startPoint, l.endPoint, l1.startPoint) {
		return true
	} else if d4 == 0 && onSegment(l.startPoint, l.endPoint, l1.endPoint) {
		return true
	} else {
		return false
	}
}

func (l0 *LineSegment) IntersectionPoint(l1 *LineSegment) util.FVec {
	if !l0.Intersects(l1) {
		return nil
	}
	var s, t float64
	x0 := (l0.endPoint.X() - l0.startPoint.X())   //difference between l0 start and end X
	x1 := (l1.endPoint.X() - l1.startPoint.X())   //difference between l1 start and end X
	xi := (l1.startPoint.X() - l0.startPoint.X()) // difference between start point X

	y0 := (l0.endPoint.Y() - l0.startPoint.Y())   //difference between l0 start and end X
	y1 := (l1.endPoint.Y() - l1.startPoint.Y())   //difference between l1 start and end X
	yi := (l1.startPoint.Y() - l0.startPoint.Y()) // difference between start point X

	//s(x0) - t(x1) = xi
	//s(y0) - t(y1) = yi
	//s = (t(y1) + yi) / y0

	tTop := (yi * x0) - (y0 * xi)
	tBottom := (x1 * y0) - (y1 * x0)
	t = tTop / tBottom
	s = ((t * y1) + yi) / y0

	return util.NewFVec2(l0.startPoint.X()+x0*s, l0.startPoint.Y()+y0*s)
}

func direction(p0, p1, p2 util.FVec) float64 {
	return util.Vec2CrossProduct(util.NewFVec2(p2.X()-p0.X(), p2.Y()-p0.Y()),
		util.NewFVec2(p1.X()-p0.X(), p1.Y()-p0.Y())) //p1 - p0
}

func onSegment(p0, p1, p2 util.FVec) bool {
	if (math.Min(p0.X(), p1.X()) <= p2.X() && p2.X() <= math.Max(p0.X(), p1.X())) &&
		(math.Min(p0.Y(), p1.Y()) <= p2.Y() && p2.Y() <= math.Max(p0.Y(), p1.Y())) {
		return true
	} else {
		return false
	}
}

func insideSegment(p0, p1, p2 util.FVec) bool {
	if p2 == p0 || p2 == p1 { //not inside, one of the endpoints
		return false
	} else {
		return onSegment(p0, p1, p2)
	}
}

func (l *LineSegment) String() string {
	return fmt.Sprintf("s: %v, e: %v",
		l.startPoint,
		l.endPoint)
}

var (
	negativeInfinityLine = NewLineSegment(
		util.NewFVec2(math.Inf(-1), math.Inf(-1)),
		util.NewFVec2(math.Inf(1), math.Inf(-1)),
	)
	positiveInfinityLine = NewLineSegment(
		util.NewFVec2(math.Inf(-1), math.Inf(1)),
		util.NewFVec2(math.Inf(1), math.Inf(1)),
	)
)

var sweeplineXPoint float64 = 0

func LineSegmentIntersection(lines []*LineSegment) *LSIState {
	lsi := NewLSIState(lines)
	for lsi.eventQueue.Length() > 0 {
		lsiIterate(lsi)
	}

	return lsi
}

func lsiIterate(lsi *LSIState) error {
	currentNode := lsi.eventQueue.Root().Minimum()
	nextPointEntry := currentNode.Key().(*eventPointEntry)
	lsi.nextPoint = nextPointEntry.FVec

	if lsi.nextPoint.X() > lsi.maxSweeplineX {
		return fmt.Errorf("reached max X")
	}

	lsi.eventQueue.Delete(currentNode)
	delete(lsi.futurePoints, nextPointEntry)
	lsi.previousPoints = append(lsi.previousPoints, lsi.currentPoint)
	lsi.currentPoint = lsi.nextPoint
	lsi.currentSweepline = lsi.currentPoint.X()

	if len(lsi.startPoints[lsi.currentPoint])+len(lsi.endPoints[lsi.currentPoint])+len(lsi.containedPoints[lsi.currentPoint]) > 1 {
		//more than one line segment starts, ends or is on this point therefore this is an intersection point
		lsi.intersectionPoints[lsi.currentPoint] = append(append(lsi.startPoints[lsi.currentPoint], lsi.endPoints[lsi.currentPoint]...),
			lsi.containedPoints[lsi.currentPoint]...) //smash them all together yeet
	}

	for _, line := range append(lsi.endPoints[lsi.currentPoint], lsi.containedPoints[lsi.currentPoint]...) {
		//Remove all the lines that are ending or are contained in this point back to line status
		lsi.lineStatus.Delete(lsi.lineStatusNodes[line])
		delete(lsi.lineStatusNodes, line)
	}

	for _, line := range append(lsi.startPoints[lsi.currentPoint], lsi.containedPoints[lsi.currentPoint]...) {
		//add all the lines that are starting or have a point in this point back to the line status
		lsi.lineStatusNodes[line] = lsi.lineStatus.Insert(newLineStatusEntry(line, lsi))
	}

	// nothing starts or continues through this point, so now we check if there are new events from the above and below lines
	// that have now become neighbours
	if len(lsi.startPoints[lsi.currentPoint])+len(lsi.containedPoints[lsi.currentPoint]) == 0 {
		cpNode := lsi.lineStatus.Insert(newEventPointEntry(lsi.currentPoint, lsi))
		anNode := cpNode.Successor()
		aboveNeighbor, ok := anNode.Key().(*lineStatusEntry)
		if ok {
			bnNode := cpNode.Predecessor()
			belowNeighbor, ok := bnNode.Key().(*lineStatusEntry)
			if ok {
				findEvent(lsi, belowNeighbor.lineSegment, aboveNeighbor.lineSegment, lsi.currentPoint)
			} else {
				panic("")
			}
		} else {
			panic("")
		}
		lsi.lineStatus.Delete(cpNode)
	} else {
		var bottomSeg, topSeg *LineSegment

		for _, lineSegment := range append(lsi.startPoints[lsi.currentPoint], lsi.containedPoints[lsi.currentPoint]...) {
			if bottomSeg == nil || compareLines(lineSegment, bottomSeg, lsi.currentSweepline+compPointFuzz) == util.ComparableLess {
				bottomSeg = lineSegment
			}
			if topSeg == nil || compareLines(lineSegment, topSeg, lsi.currentSweepline+compPointFuzz) == util.ComparableGreater {
				topSeg = lineSegment
			}
		}

		bsp := lsi.lineStatusNodes[bottomSeg].Predecessor()
		if bsp != nil && !bsp.Nil() {
			findEvent(lsi, bsp.Key().(*lineStatusEntry).lineSegment,
				bottomSeg, lsi.currentPoint)
		}

		tsp := lsi.lineStatusNodes[topSeg].Successor()
		if tsp != nil && !tsp.Nil() {
			findEvent(lsi, tsp.Key().(*lineStatusEntry).lineSegment,
				topSeg, lsi.currentPoint)
		}
	}
	return nil
}

type LSIState struct {
	currentSweepline   float64
	maxSweeplineX      float64
	lines              []*LineSegment
	intersectionPoints PointSegments
	lineStatus         tree.BinarySearchTree
	startPoints        PointSegments
	endPoints          PointSegments
	containedPoints    PointSegments
	lineStatusNodes    map[*LineSegment]tree.BSTNode
	eventQueue         tree.BinarySearchTree
	previousPoints     []util.FVec
	currentPoint       util.FVec
	nextPoint          util.FVec
	futurePoints       map[*eventPointEntry]util.FVec
}

func NewLSIState(lines []*LineSegment) *LSIState {
	s := &LSIState{
		lines: lines,
	}
	s.maxSweeplineX = math.Inf(1)
	s.Clear()
	return s
}

func (s *LSIState) Clear() {
	s.currentSweepline = 0
	s.intersectionPoints = make(PointSegments)
	s.initLineStatus()
	s.lineStatusNodes = make(map[*LineSegment]tree.BSTNode)
	s.containedPoints = make(PointSegments)
	s.previousPoints = make([]util.FVec, 0, 10)
	s.currentPoint = nil
	s.nextPoint = nil
	s.futurePoints = make(map[*eventPointEntry]util.FVec)
	s.initEventQueue()
}

func (s *LSIState) Lines() []*LineSegment {
	return s.lines
}

func (s *LSIState) SetMaxSweeplineX(x float64) {
	s.maxSweeplineX = x
}

func (s *LSIState) CurrentSweepline() float64 {
	return s.currentSweepline
}

func (s *LSIState) PreviousPoints() []util.FVec {
	return s.previousPoints
}

func (s *LSIState) CurrentPoint() util.FVec {
	return s.currentPoint
}

func (s *LSIState) NextPoint() util.FVec {
	return s.nextPoint
}

func (s *LSIState) FuturePoints() []util.FVec {
	points := make([]util.FVec, len(s.futurePoints))
	i := 0
	for _, p := range s.futurePoints {
		points[i] = p
		i++
	}
	return points
}

func (s *LSIState) IntersectionPoints() PointSegments {
	return s.intersectionPoints
}

func (s *LSIState) ActiveLines() []*LineSegment {
	al := make([]*LineSegment, len(s.lineStatusNodes))
	i := 0
	for line := range s.lineStatusNodes {
		if line != negativeInfinityLine && line != positiveInfinityLine {
			al[i] = line
			i++
		}
	}

	return al
}

func (s *LSIState) Run() {
	for s.eventQueue.Length() > 0 {
		if err := lsiIterate(s); err != nil {
			break
		}
	}
}

func (s *LSIState) initEventQueue() {
	s.eventQueue = tree.NewRBTree()
	s.startPoints = make(PointSegments)
	s.endPoints = make(PointSegments)
	for _, line := range s.lines {
		addEventPoint(s, line.Start())
		s.startPoints[line.Start()] = append(s.startPoints[line.Start()], line)
		addEventPoint(s, line.End())
		s.endPoints[line.End()] = append(s.endPoints[line.End()], line)
	}
}

func (s *LSIState) initLineStatus() {
	s.lineStatus = tree.NewRBTree()

	nlineEntry := newLineStatusEntry(negativeInfinityLine, s)
	s.lineStatus.Insert(nlineEntry)
	plineEntry := newLineStatusEntry(positiveInfinityLine, s)
	s.lineStatus.Insert(plineEntry)
}

type lineStatusEntry struct {
	lineSegment *LineSegment
	state       SweeplineState
}

func newLineStatusEntry(segment *LineSegment, state *LSIState) *lineStatusEntry {
	lse := lineStatusEntry{
		lineSegment: segment,
		state:       state,
	}
	return &lse
}

func (le *lineStatusEntry) String() string {
	return le.lineSegment.String()
}

func lineYPoint(l *LineSegment, x float64) (y util.ComparableFloat) {
	v := util.ComparableFloat((l.gradient * x) + l.yIntercept)
	return v
}

func (le *lineStatusEntry) Compare(comp util.Comparable) util.ComparableResult {
	switch le1 := comp.(type) {
	case *lineStatusEntry:
		//the fuzz is added so that the line is slightly in front of where the sweepline is, to break ties at intersection points
		return compareLines(le.lineSegment, le1.lineSegment, le.state.CurrentSweepline()+compPointFuzz)
	case *eventPointEntry:
		return lineYPoint(le.lineSegment, le.state.CurrentSweepline()).Compare(util.ComparableFloat(le1.Y()))
	default:
		panic("invalid compare")
	}
}

var compPointFuzz float64 = 0.00001

func comparePoints(p0, p1 util.FVec) util.ComparableResult {
	dx := math.Abs(p1.X() - p0.X())
	dy := math.Abs(p1.X() - p0.X())
	if p0.X() == p1.X() && p0.Y() == p1.Y() || (dx < compPointFuzz && dy < compPointFuzz) {
		return util.ComparableEqual
	} else if p0.X() < p1.X() {
		return util.ComparableLess
	} else if p0.X() == p1.X() {
		if p0.Y() < p1.Y() {
			return util.ComparableLess
		} else {
			return util.ComparableGreater
		}
	} else {
		return util.ComparableGreater
	}
}

func compareLines(l0, l1 *LineSegment, xPoint float64) util.ComparableResult {
	return lineYPoint(l0, xPoint).Compare(lineYPoint(l1, xPoint))
}

// Initialise line status with lines at minus infinity and plus infinity to act as sentinels
type eventPointEntry struct {
	util.FVec
	state SweeplineState
}

func newEventPointEntry(point util.FVec, state SweeplineState) *eventPointEntry {
	e := eventPointEntry{
		FVec:  point,
		state: state,
	}
	return &e
}

func (epe *eventPointEntry) String() string {
	return fmt.Sprint(epe.FVec)
}

func (ep *eventPointEntry) Compare(p util.Comparable) util.ComparableResult {
	switch ep1 := p.(type) {
	case *eventPointEntry:
		return comparePoints(ep, ep1)
	case *lineStatusEntry:
		return util.ComparableFloat(ep.Y()).Compare(lineYPoint(ep1.lineSegment, ep.state.CurrentSweepline()))
	default:
		panic(fmt.Sprintf("invalid compare with %v\n", p))
	}
}

type PointSegments map[util.FVec][]*LineSegment

func addEventPoint(lsi *LSIState, eventPoint util.FVec) {
	epe := newEventPointEntry(eventPoint, lsi)
	exists, _ := tree.BSTNodeSearch(lsi.eventQueue.Root(), epe)
	if exists == nil || exists.Nil() {
		lsi.eventQueue.Insert(epe)
		lsi.futurePoints[epe] = eventPoint
	}
}

func findEvent(lsi *LSIState, seg0, seg1 *LineSegment, point util.FVec) {
	if !seg0.Intersects(seg1) {
		return
	}
	intersectionPoint := seg0.IntersectionPoint(seg1)
	if (intersectionPoint.X() > point.X()) ||
		(intersectionPoint.X() == point.X() && intersectionPoint.Y() > point.Y()) {
		//if the intersection point is in front of the line
		//or if it's on the sweepline but above our currentpoint
		addEventPoint(lsi, intersectionPoint) //this won't add the point if it's already there
	}
	if insideSegment(seg0.startPoint, seg0.endPoint, intersectionPoint) {
		found := false
		for _, seg := range lsi.containedPoints[intersectionPoint] {
			if seg0 == seg {
				found = true
			}
		}
		if !found {
			lsi.containedPoints[intersectionPoint] = append(lsi.containedPoints[intersectionPoint], seg0)
		}
	}
	if insideSegment(seg1.startPoint, seg1.endPoint, intersectionPoint) {
		found := false
		for _, seg := range lsi.containedPoints[intersectionPoint] {
			if seg1 == seg {
				found = true
			}
		}
		if !found {
			lsi.containedPoints[intersectionPoint] = append(lsi.containedPoints[intersectionPoint], seg1)
		}
	}
}
