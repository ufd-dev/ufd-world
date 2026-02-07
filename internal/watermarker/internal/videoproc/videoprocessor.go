package videoproc

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

const bufferSize = 100
const jobTimeout = 3 * time.Minute

// job defines the input parameters for FFmpeg
type job struct {
	source  string
	overlay string
	target  string
}

// videoProcessor manages the job queue and worker
type videoProcessor struct {
	queue chan job
}

var vpSingleton *videoProcessor

// GetSingleton gets a singleton videoProcessor and inits a new one as necessary
func GetSingleton() *videoProcessor {
	if vpSingleton == nil {
		vpSingleton = &videoProcessor{
			queue: make(chan job, bufferSize),
		}
		// Start the single worker goroutine
		go vpSingleton.startWorker()
	}
	return vpSingleton
}

// AddJob is thread-safe and can be called from any HTTP handler
func (vp *videoProcessor) AddJob(source, overlay, target string) error {
	job := job{source: source, overlay: overlay, target: target}

	// This will not block the caller unless the buffer is completely full
	select {
	case vp.queue <- job:
		return nil
	default:
		return fmt.Errorf("Queue full! Could not add job: %s", source)
	}
}

func (vp *videoProcessor) startWorker() {
	// This loop runs for the life of the application
	for job := range vp.queue {
		vp.runFFmpeg(job)
	}
}

func (vp *videoProcessor) runFFmpeg(job job) {
	// Context ensures the job is killed if it hangs over 3 minutes
	ctx, cancel := context.WithTimeout(context.Background(), jobTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", job.source,
		"-i", job.overlay,
		"-filter_complex", "overlay=0:0:format=auto,format=yuv420p",
		"-movflags", "faststart",
		"-y",
		job.target,
	)

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("TIMEOUT: FFmpeg killed for %s after 3m", job.source)
		} else {
			log.Printf("ERROR: FFmpeg failed for %s: %v", job.source, err)
		}
		return
	}

	log.Printf("SUCCESS: Created %s", job.target)
}
