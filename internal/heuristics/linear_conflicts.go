package heuristics

func linearConflicts(grid []int, size int, depth int) float32 {
	conflicts := 0
	for x1 := 0; x1 < size; x1++ {
		for y1 := 0; y1 < size-1; y1++ {
			if grid[get1d(x1, y1, size)] != 0 {
				tmp := finalPos[grid[get1d(x1, y1, size)]]
				targetx, targety := tmp[0], tmp[1]
				if (x1 == targetx) != (y1 == targety) {
					if (x1 == targetx) && (y1 < targety) {
						// Case 1: my x is ok, my target in on my right
						// I check the other ones in this line
						for j := y1; j < size; j++ {
							currPos := get1d(x1, j, size)
							if grid[currPos] != 0 {
								currGoalPos := finalPos[grid[currPos]]
								// if his target is in this line and before mine
								if currGoalPos[0] == x1 && currGoalPos[1] < targety {
									conflicts++
								}
							}
						}
					}
					if (y1 == targety) && (x1 < targetx) {
						// Case 2: my y is ok, my target is under me
						// I check the other ones in this col
						for i := x1; i < size; i++ {
							currPos := get1d(i, y1, size)
							if grid[currPos] != 0 {
								currGoalPos := finalPos[grid[currPos]]
								// if his target is in this col and before mine
								if currGoalPos[1] == y1 && currGoalPos[0] <= targetx {
									conflicts++
								}
							}
						}
					}
				}
			}
		}
	}
	return manhattan(grid, size, depth) + 2*float32(conflicts)
}

func linearConflictsA(grid []int, size int, depth int) float32 {
	return linearConflicts(grid, size, depth) + float32(depth)
}
