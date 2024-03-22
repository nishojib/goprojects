package solver

import (
	"image"
	"log"
	"sync"
)

// explore one path and publish to the s.pathsToExplore channel
// any branch we discover that we don't take.
func (s *Solver) explore(pathToBranch *path) {
	if pathToBranch == nil {
		// This is a safety net. It should be used, but when it's needed, at least it's there.
		return
	}

	currentPosition := pathToBranch.at

	for {
		// Paint the current pixel as explored.
		s.maze.Set(currentPosition.X, currentPosition.Y, s.palette.explored)

		// Let's first check whether we should quit.
		select {
		case <-s.quit:
			return
		case s.exploredPixels <- currentPosition:
			// continue the exploration
		}

		// We know we'll have up to 3 new neighbors to explore.
		candidates := make([]image.Point, 0, 3)

		for _, neighbor := range neighbors(currentPosition) {
			if pathToBranch.isPreviousStep(neighbor) {
				// Let's not return to the previous position.
				continue
			}

			// Look at the color of this pixel.
			// RGBAAt returns a color.RGBA{} zero value if the pixel is outside the bounds of the image.
			switch s.maze.RGBAAt(neighbor.X, neighbor.Y) {
			case s.palette.treasure:
				s.mutex.Lock()
				defer s.mutex.Unlock()
				if s.solution == nil {
					s.solution = &path{previousStep: pathToBranch, at: neighbor}
					log.Printf("Treasure found at %v!", neighbor)
					close(s.quit)
				}
				return
			case s.palette.path:
				candidates = append(candidates, neighbor)
			}
		}

		if len(candidates) == 0 {
			log.Printf("I must have taken the wrong turn at position %v", currentPosition)
			return
		}

		for _, candidate := range candidates[1:] {
			branch := &path{previousStep: pathToBranch, at: candidate}
			// We are sure we send to pathsToExplore only when the quit channel isn't closed.
			// A goroutine might have found the treasure since the check at the start of the loop.
			select {
			case <-s.quit:
				log.Printf(
					"I am an unlucky branch, someone else found the treasure, I give up at position %v.",
					currentPosition,
				)
				return
			case s.pathsToExplore <- branch:
				// continue execution after the select block.
			}
		}

		pathToBranch = &path{previousStep: pathToBranch, at: candidates[0]}
		currentPosition = candidates[0]
	}
}

// listenToBranches creates a new goroutine for each branch published in s.pathsToExplore.
func (s *Solver) listenToBranches() {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for {
		select {
		// s.quit will never return a value, unless something writes in it (which we don't do)
		// or it has been closed, which we do when we find the treasure.
		case <-s.quit:
			log.Println("there treasure has been found, stopping the worker")
			return
		case p := <-s.pathsToExplore:
			wg.Add(1)
			go func(path *path) {
				defer wg.Done()
				s.explore(path)
			}(p)
		}
	}
}
