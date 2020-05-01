/**
 * @ Author: Loi Nguyen
 * @ Create Time: 2020-04-29 20:31:00
 * @ Modified by: Loi Nguyen
 * @ Modified time: 2020-04-30 17:03:45
 * @ Description:
 */

package main

import (
	"fmt"
	"os"

	"github.com/golang-collections/collections/queue"
)

type ii struct {
	x int
	y int
}

var mapName string
var nRow, nCol, nGaz, xStart, yStart int
var matrix [][]int
var gazList []ii
var gazEdge [][]int
var gazStar [][]ii
var mapGazID [][]int
var mapVisited [][]bool
var bfsVisited [][]bool
var starCollected [][]bool
var forDirections []ii
var backtrack [][]ii
var depth [][]int
var gazOrder []int
var gazVisited []bool
var cntGazVisited int
var hasEdge [][]bool
var starPath [][][]ii
var gazPath [][][]ii

func main() {
	mapName = "map.txt"
	
	input()
	
	buildGazMap()
	
	path1 := toNearestGaz()
	nearestGazPos := path1[len(path1)-1]
	nearestGazID := mapGazID[nearestGazPos.x][nearestGazPos.y]

	dfsGaz(nearestGazID)

	path2 := genPathByGazOrder()

	path := append(path1, path2...)

	output(path)
}

func genPathByGazOrder() []ii {
	path := []ii{}

	for i := 0; i < len(gazOrder); i++ {
		// Collect star arround gazOrder[i]
		gazID := gazOrder[i]
		for j := 0; j < len(gazStar[gazID]); j++ {
			path = append(path, starPath[gazID][j]...)

			backPath := reversePath(gazList[gazID], starPath[gazID][j])

			path = append(path, backPath...)
		}

		// Go to next gaz station
		if i != len(gazOrder)-1 {
			from := gazOrder[i]
			to := gazOrder[i+1]

			path = append(path, gazPath[from][to]...)
		}
	}

	return path
}

func dfsGaz(gazID int) {
	if gazVisited[gazID] == true {
		return
	}

	gazVisited[gazID] = true
	cntGazVisited++
	gazOrder = append(gazOrder, gazID)

	for _, e := range gazEdge[gazID] {
		dfsGaz(e)

		if cntGazVisited != len(gazList) {
			gazOrder = append(gazOrder, gazID)
		}
	}
}

func toNearestGaz() []ii {
	clearMapVisited()
	q := queue.New()
	q.Enqueue(ii{xStart, yStart})

	for q.Len() != 0 {
		node := q.Dequeue().(ii)
		x := node.x
		y := node.y

		if x < 0 || x >= nRow || y < 0 || y >= nCol ||
			mapVisited[x][y] || matrix[x][y] == 0 {
			continue
		}

		// Nearest gaz
		if matrix[x][y] == 2 {
			return tracePath(xStart, yStart, x, y)
		}

		mapVisited[x][y] = true

		for i := 0; i < 4; i++ {
			newX := x + forDirections[i].x
			newY := y + forDirections[i].y

			if newX < 0 || newX >= nRow || newY < 0 || newY >= nCol ||
				mapVisited[newX][newY] || matrix[newX][newY] == 0 {
				continue
			}

			q.Enqueue(ii{newX, newY})
			backtrack[newX][newY] = ii{x, y}
		}
	}

	return nil
}

func buildGazMap() {
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if matrix[i][j] == 2 {
				gazList = append(gazList, ii{i, j})
				mapGazID[i][j] = len(gazList) - 1
			}
		}
	}

	gazEdge = make([][]int, len(gazList))
	gazStar = make([][]ii, len(gazList))
	starPath = make([][][]ii, len(gazList))
	gazPath = make([][][]ii, len(gazList))
	gazVisited = make([]bool, len(gazList))
	hasEdge = make([][]bool, len(gazList))

	for i := 0; i < len(gazList); i++ {
		gazEdge[i] = make([]int, 0)
		gazStar[i] = make([]ii, 0)
		starPath[i] = make([][]ii, 0)
		gazPath[i] = make([][]ii, len(gazList))
		hasEdge[i] = make([]bool, len(gazList))
	}

	for i := 0; i < len(gazList); i++ {
		bfsGaz(i, gazList[i].x, gazList[i].y)
	}
}

func clearMapVisited() {
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			mapVisited[i][j] = false
			bfsVisited[i][j] = false
		}
	}
}

func reverseSlice(a []ii) []ii {
	b := make([]ii, len(a))
	copy(b, a)

	// Reverse path
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return b
}

func reversePath(src ii, path []ii) []ii {
	removedDestPath := path[:len(path)-1]
	newPath := reverseSlice(removedDestPath)
	newPath = append(newPath, src)

	return newPath
}

func tracePath(fromX, fromY, toX, toY int) []ii {
	path := []ii{{toX, toY}}

	for back := backtrack[toX][toY]; back.x != fromX || back.y != fromY; back = backtrack[toX][toY] {
		path = append(path, back)
		toX = back.x
		toY = back.y
	}

	path = reverseSlice(path)
	return path
}

func bfsGaz(gazID, x0, y0 int) {
	clearMapVisited()
	q := queue.New()
	q.Enqueue(ii{x0, y0})
	depth[x0][y0] = 0
	gazPos := gazList[gazID]

	for q.Len() != 0 {
		node := q.Dequeue().(ii)
		x := node.x
		y := node.y
		if x < 0 || x >= nRow || y < 0 || y >= nCol ||
			bfsVisited[x][y] || depth[x][y] > nGaz || matrix[x][y] == 0 {
			continue
		}

		bfsVisited[x][y] = true
		if !(x == x0 && y == y0) {
			switch matrix[x][y] {
			case 2: // Gaz station
				gazID2 := mapGazID[x][y]
				if hasEdge[gazID][gazID2] {
					continue
				}

				hasEdge[gazID][gazID2] = true
				hasEdge[gazID2][gazID] = true

				gazEdge[gazID] = append(gazEdge[gazID], gazID2)
				gazEdge[gazID2] = append(gazEdge[gazID2], gazID)
				path := tracePath(gazPos.x, gazPos.y, x, y)
				gazPath[gazID][gazID2] = path
				gazPath[gazID2][gazID] = reversePath(gazPos, path)
			case 3: // Star
				if starCollected[x][y] == false && depth[x][y] <= nGaz/2 {
					gazStar[gazID] = append(gazStar[gazID], ii{x, y})
					starCollected[x][y] = true

					// Trace path
					starPath[gazID] = append(starPath[gazID], make([]ii, 0))
					starPath[gazID][len(starPath[gazID])-1] = tracePath(gazPos.x, gazPos.y, x, y)
				}
			}
		}

		for i := 0; i < 4; i++ {
			newX := x + forDirections[i].x
			newY := y + forDirections[i].y

			if newX < 0 || newX >= nRow || newY < 0 || newY >= nCol ||
				bfsVisited[newX][newY] || matrix[x][y] == 0 {
				continue
			}

			depth[newX][newY] = depth[x][y] + 1
			backtrack[newX][newY] = ii{x, y}

			q.Enqueue(ii{newX, newY})
		}
	}
}

func input() {
	gazOrder = []int{}
	forDirections = []ii{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
	os.Stdin, _ = os.OpenFile(mapName, os.O_RDONLY|os.O_CREATE, 0666)
	fmt.Println("")
	fmt.Scanf("%d %d\n%d\n%d %d", &xStart, &yStart, &nGaz, &nRow, &nCol)

	xStart--
	yStart--

	for i := 0; i < nRow; i++ {
		matrix = append(matrix, make([]int, nCol))
		mapVisited = append(mapVisited, make([]bool, nCol))
		mapGazID = append(mapGazID, make([]int, nCol))
		depth = append(depth, make([]int, nCol))
		backtrack = append(backtrack, make([]ii, nCol))
		starCollected = append(starCollected, make([]bool, nCol))
		bfsVisited = append(bfsVisited, make([]bool, nCol))
		for j := 0; j < nCol; j++ {
			fmt.Scanf("%d", &matrix[i][j])
			mapVisited[i][j] = false
			bfsVisited[i][j] = false
		}
	}
}

func output(path []ii) {
	os.Stdout, _ = os.OpenFile("result_"+mapName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	for i := 0; i < len(path); i++ {
		fmt.Printf("%d %d\n", path[i].x+1, path[i].y+1)
	}
}
