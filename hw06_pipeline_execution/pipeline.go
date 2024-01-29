package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stageExecute(done, stage, in)
	}

	return in
}

func stageExecute(done In, stage Stage, in In) Out {
	resultChannel := make(Bi)
	outChannel := stage(in)

	go func() {
		defer close(resultChannel)

		for {
			select {
			case <-done:
				return
			case value, ok := <-outChannel:
				if !ok {
					return
				}
				resultChannel <- value
			}
		}
	}()

	return resultChannel
}
