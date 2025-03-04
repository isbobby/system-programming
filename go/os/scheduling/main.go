package main

import "context"

func main() {
	newProcessor := New()
	ctx := context.Background()
	newProcessor.RunFifo(ctx, caseThree)
}
