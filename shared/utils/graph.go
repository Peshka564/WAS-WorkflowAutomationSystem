package utils

import "fmt"

// TODO: This only works for one single listener in the whole worflow
func TopologicalSort(graph map[string][]string) ([]string, error) {
	fmt.Println(graph)
    inDegree := make(map[string]int)
    
    for node := range graph {
        inDegree[node] = 0
    }
    
    for _, neighbors := range graph {
        for _, neighbor := range neighbors {
            inDegree[neighbor]++
        }
    }

    var queue []string
    for node, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, node)
        }
    }

    var sortedOrder []string
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:] // Dequeue

        sortedOrder = append(sortedOrder, current)

        for _, neighbor := range graph[current] {
            inDegree[neighbor]--
            
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    // Cycle Detection Check
    if len(sortedOrder) != len(inDegree) {
        return nil, fmt.Errorf("cycle detected or graph is invalid")
    }

    return sortedOrder, nil
}