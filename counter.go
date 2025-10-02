package main

import (
	"fmt"
	"os"
	"strconv"

	"gioui.org/widget"
)

type Counter struct {
	fileName        *string
	incrementButton *widget.Clickable
	decrementButton *widget.Clickable
	value           int
}

func (counter *Counter) increment() error {
	data := []byte(strconv.FormatInt(int64(counter.value+1), 10))
	name := fmt.Sprintf("%s/%s", folderPath, *counter.fileName)

	w, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed open file %s: %w", name, err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed write file %s: %w", name, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed close file %s: %w", name, err)
	}

	counter.value += 1

	return nil
}

func (counter *Counter) decrement() error {
	data := []byte(strconv.FormatInt(int64(counter.value-1), 10))
	name := fmt.Sprintf("%s/%s", folderPath, *counter.fileName)

	w, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed open file %s: %w", name, err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed write file %s: %w", name, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed close file %s: %w", name, err)
	}

	counter.value -= 1

	return nil
}
