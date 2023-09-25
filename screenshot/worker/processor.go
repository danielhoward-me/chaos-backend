package worker

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/utils"

	"context"
	"fmt"
	"os"

	"github.com/chromedp/chromedp"
)

func run() {
	workerRunning = true

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(1920, 1080),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	for {
		if len(jobs) == 0 {
			break
		}

		var job Job
		job, jobs = jobs[0], jobs[1:]

		fmt.Printf("Running screenshot job for %s\n", job.hash)
		if err := process(job, ctx); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Finished screenshot job for %s\n", job.hash)
	}

	workerRunning = false
}

func process(job Job, ctx context.Context) error {
	url := getUrl()

	var buf []byte

	tasks := getTasks(url, job.data, &buf)
	if err := chromedp.Run(ctx, tasks); err != nil {
		return fmt.Errorf("failed to run screenshot tasks for %s: %s", job.hash, err)
	}

	path := utils.Path(job.hash)
	if err := os.WriteFile(path, buf, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write screenshot file for %s: %s", job.hash, err)
	}

	return nil
}

func getUrl() string {
	url := "https://chaos.danielhoward.me"
	if chaosDevPort != 0 {
		url = fmt.Sprintf("http://local.danielhoward.me:%d", chaosDevPort)
	}

	return fmt.Sprintf("%s/?screenshot-worker", url)
}

func getTasks(url string, data string, buf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Evaluate(fmt.Sprintf("window.prepareScreenshot(`%s`)", data), nil),
		chromedp.FullScreenshot(buf, 99),
	}
}
