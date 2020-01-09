/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: file.go
 * @time: 2020/1/7 15:52
 */
package logadapter

import (
	`bytes`
	`errors`
	`fmt`
	`io`
	`os`
	`path`
	`path/filepath`
	`strings`
	`sync`
	`time`
)

type FileWriter struct {
	sync.RWMutex

	fileWriter       *os.File
	Filename         string
	suffix           string
	fileNameOnly     string
	maxSizeCurSize   int
	dailyOpenTime    time.Time
	dailyOpenDate    int
	hourlyOpenTime   time.Time
	hourlyOpenDate   int
	maxLinesCurLines int
	MaxLines         int
	MaxSize          int
	MaxFilesCurFiles int
	MaxFiles         int
	MaxHours         int
	MaxDays          int
	// Hourly           bool
	Daily            bool
	Rotate bool
}

func (p *FileWriter) Close() error {
	return p.fileWriter.Close()
}

func (p *FileWriter) Write(msg []byte) (n int, err error) {
	// hd, d, h := formatTimeHeader(when)
	// msg = string(hd) + msg + "\n"
	now := time.Now()
	d := now.Day()
	// h := now.Hour()
	if p.Rotate {
		p.RLock()
		/*if p.needRotateHourly(len(msg), h) {
			p.RUnlock()
			p.Lock()
			if p.needRotateHourly(len(msg), h) {
				if err := p.doRotate(now); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", p.Filename, err)
				}
			}
			p.Unlock()
		} else*/ if p.needRotateDaily(len(msg), d) {
			p.RUnlock()
			p.Lock()
			if p.needRotateDaily(len(msg), d) {
				if err := p.doRotate(now); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", p.Filename, err)
				}
			}
			p.Unlock()
		} else {
			p.RUnlock()
		}
	}

	p.Lock()
	_, err = p.fileWriter.Write([]byte(msg))
	if err == nil {
		p.maxLinesCurLines++
		p.maxSizeCurSize += len(msg)
	}
	p.Unlock()
	return len(msg), err
}

func (p *FileWriter) Init(filename string) error {
	if len(filename) == 0 {
		return errors.New("must have filename")
	}
	p.Filename = filename
	p.suffix = filepath.Ext(filename)
	p.fileNameOnly = strings.TrimSuffix(filename, p.suffix)
	if p.suffix == "" {
		p.suffix = ".log"
	}
	err := p.startLogger()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (p *FileWriter) startLogger() error {
	file, err := p.createLogFile()
	if err != nil {
		return err
	}
	if p.fileWriter != nil {
		p.fileWriter.Close()
	}
	p.fileWriter = file
	return p.initFd()
}

func (p *FileWriter) initFd() error {
	fd := p.fileWriter
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}
	p.Daily = true
	p.MaxLines = 1024000000
	p.MaxSize = 1024000000
	p.MaxDays = 3
	p.maxSizeCurSize = int(fInfo.Size())
	p.dailyOpenTime = time.Now()
	p.dailyOpenDate = p.dailyOpenTime.Day()
	p.hourlyOpenTime = time.Now()
	p.hourlyOpenDate = p.hourlyOpenTime.Hour()
	p.maxLinesCurLines = 0
	/*if p.Hourly {
		go p.hourlyRotate(p.hourlyOpenTime)
	} else */if p.Daily {
		go p.dailyRotate(p.dailyOpenTime)
	}
	if fInfo.Size() > 0 && p.MaxLines > 0 {
		count, err := p.lines()
		if err != nil {
			return err
		}
		p.maxLinesCurLines = count
	}
	return nil
}

func (p *FileWriter) lines() (int, error) {
	fd, err := os.Open(p.Filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func (p *FileWriter) createLogFile() (*os.File, error) {
	pt := path.Dir(p.Filename)
	os.MkdirAll(pt, os.ModePerm)

	fd, err := os.OpenFile(p.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err == nil {
		os.Chmod(p.Filename, os.ModePerm)
	}
	return fd, err
}

func (p *FileWriter) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	p.Lock()
	if p.needRotateDaily(0, time.Now().Day()) {
		if err := p.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", p.Filename, err)
		}
	}
	p.Unlock()
}
/*
func (p *FileWriter) hourlyRotate(openTime time.Time) {
	y, m, d := openTime.Add(1 * time.Hour).Date()
	h, _, _ := openTime.Add(1 * time.Hour).Clock()
	nextHour := time.Date(y, m, d, h, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextHour.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	p.Lock()
	if p.needRotateHourly(0, time.Now().Hour()) {
		if err := p.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", p.Filename, err)
		}
	}
	p.Unlock()
}*/

func (p *FileWriter) needRotateDaily(size int, day int) bool {
	return (p.MaxLines > 0 && p.maxLinesCurLines >= p.MaxLines) ||
		(p.MaxSize > 0 && p.maxSizeCurSize >= p.MaxSize) ||
		(p.Daily && day != p.dailyOpenDate)
}

// func (p *FileWriter) needRotateHourly(size int, hour int) bool {
// 	return (p.MaxLines > 0 && p.maxLinesCurLines >= p.MaxLines) ||
// 		(p.MaxSize > 0 && p.maxSizeCurSize >= p.MaxSize) ||
// 		(p.Hourly && hour != p.hourlyOpenDate)
// }

// DoRotate means it need to write file in new file.
// new file name like xx.2013-01-01.log (daily) or xx.001.log (by line or size)
func (p *FileWriter) doRotate(logTime time.Time) error {
	// file exists
	// Find the next available number
	num := p.MaxFilesCurFiles + 1
	fName := ""
	format := ""
	var openTime time.Time
	// rotatePerm, err := strconv.ParseInt(w.RotatePerm, 8, 64)
	// if err != nil {
	// 	return err
	// }

	_, err := os.Lstat(p.Filename)
	if err != nil {
		// even if the file is not exist or other ,we should RESTART the logger
		goto RESTART_LOGGER
	}

	/*if p.Hourly {
		format = "2006010215"
		openTime = p.hourlyOpenTime
	} else*/ if p.Daily {
		format = "2006-01-02"
		openTime = p.dailyOpenTime
	}

	// only when one of them be setted, then the file would be splited
	if p.MaxLines > 0 || p.MaxSize > 0 {
		for ; err == nil && num <= p.MaxFiles; num++ {
			fName = p.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logTime.Format(format), num, p.suffix)
			_, err = os.Lstat(fName)
		}
	} else {
		fName = p.fileNameOnly + fmt.Sprintf(".%s.%03d%s", openTime.Format(format), num, p.suffix)
		_, err = os.Lstat(fName)
		p.MaxFilesCurFiles = num
	}

	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("rotate: Cannot find free log number to rename %s", p.Filename)
	}

	// close fileWriter before rename
	p.fileWriter.Close()

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	err = os.Rename(p.Filename, fName)
	if err != nil {
		goto RESTART_LOGGER
	}

	err = os.Chmod(fName, os.ModePerm)

RESTART_LOGGER:

	startLoggerErr := p.startLogger()
	go p.deleteOldLog()

	if startLoggerErr != nil {
		return fmt.Errorf("rotate StartLogger: %s", startLoggerErr)
	}
	if err != nil {
		return fmt.Errorf("rotate: %s", err)
	}
	return nil
}

func (p *FileWriter) deleteOldLog() {
	dir := filepath.Dir(p.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil {
			return
		}
		/*if p.Hourly {
			if !info.IsDir() && info.ModTime().Add(1 * time.Hour * time.Duration(p.MaxHours)).Before(time.Now()) {
				if strings.HasPrefix(filepath.Base(path), filepath.Base(p.fileNameOnly)) &&
					strings.HasSuffix(filepath.Base(path), p.suffix) {
					os.Remove(path)
				}
			}
		} else*/ if p.Daily {
			if !info.IsDir() && info.ModTime().Add(24 * time.Hour * time.Duration(p.MaxDays)).Before(time.Now()) {
				if strings.HasPrefix(filepath.Base(path), filepath.Base(p.fileNameOnly)) &&
					strings.HasSuffix(filepath.Base(path), p.suffix) {
					os.Remove(path)
				}
			}
		}
		return
	})
}

// Destroy close the file description, close file writer.
func (p *FileWriter) Destroy() {
	p.Close()
}
