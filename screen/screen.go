// This code is largely inspired by periph.io ssd1306 package
// I just add the parameter address to NewI2C in newI2CAddress

            // Original code License \\
// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package screen

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/nfnt/resize"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
)

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
    src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
    img := image.NewGray(image.Rect(0, 0, w, h))
    r := src.Bounds()
    r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
    draw.Draw(img, r, src, image.Point{}, draw.Src)
    return img
}

type Screen struct {
    dev Dev
}

func resetDisplay() {
	// Lookup a pin by its number:
    p := gpioreg.ByName("GPIO4")
    if p == nil {
        log.Fatal("Failed to find display reset pin GPIO4")
    }
	// Set the output to low:
    if err := p.Out(gpio.Low); err != nil {
        log.Fatal(err)
    }

    // Wait a little bit before re-enabling the display
    time.Sleep(100* time.Millisecond)

	// Set the output to high:
    if err := p.Out(gpio.High); err != nil {
        log.Fatal(err)
    }
}

func NewScreen() (*Screen, error) {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

	// Reset the display
	resetDisplay()

    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")
    if err != nil {
        log.Fatal(err)
    }

    // Open a handle to a ssd1306 connected on the I²C bus:
    dev, err := newI2CAddress(bus, &DefaultOpts, 0x3D)
    if err != nil {
        log.Fatal(err)
    }
    screen := &Screen{
        dev: *dev,
    }

    return screen, nil
}

type Msg struct {
	x, y int
	msg string
}

func (msg *Msg) SetMsg(x int, y int, message string, v ...interface{}){
	msg.msg = message
	msg.x = x
	msg.y = y
}

func (screen *Screen) DisplayMessage(messages []Msg) {
    // Draw on it.
    // log.Printf("Displaying message %v", msg)
	img := image.NewGray(image.Rect(0, 0, DefaultOpts.W, DefaultOpts.H))
	
	f2 := basicfont.Face7x13
	drawer := font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{image1bit.On},
		Face: f2,
		Dot:  fixed.P(0, img.Bounds().Dy()-1-f2.Descent),
	}

	for _, msg := range messages {
		drawer.Dot.X = fixed.Int26_6(msg.x<<6)
		drawer.Dot.Y = fixed.Int26_6(msg.y<<6)
		drawer.DrawString(msg.msg)
	}
		
	if err := screen.dev.Draw(img.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}
}
    
func (screen *Screen) DisplayGif(gif_path string, quit chan bool){
    // Decodes an animated GIF as specified on the command line:
    f, err := os.Open(gif_path)
    if err != nil {
        log.Fatal(err)
    }

	log.Printf("Decoding image %v", gif_path)

    g, err := gif.DecodeAll(f)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

	log.Printf("Converting image %v", gif_path)
    // Converts every frame to image.Gray and resize them:
    imgs := make([]*image.Gray, len(g.Image))
    for i := range g.Image {
        imgs[i] = convertAndResizeAndCenter(DefaultOpts.W, DefaultOpts.H, g.Image[i])
    }

	log.Printf("Displaying image %v", gif_path)
    // Display the frames in a loop:
    for i := 0; ; i++ {
        select {
        case <-quit:
			os.Clearenv()
            return
        default:
            index := i % len(imgs)
            c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
            img := imgs[index]
            screen.dev.Draw(img.Bounds(), img, image.Point{})
            <-c
        }
    }
}

// newI2CAddress returns a Dev object that communicates over I²C to a SSD1306 display
// controller.
func newI2CAddress(i i2c.Bus, opts *Opts, address uint16) (*Dev, error) {
	// Maximum clock speed is 1/2.5µs = 400KHz.
	return newDev(&i2c.Dev{Bus: i, Addr: address}, opts, false, nil)
}


// Package ssd1306 \\
const (
	_CHARGEPUMP          = 0x8D
	_COLUMNADDR          = 0x21
	_COMSCANDEC          = 0xC8
	_COMSCANINC          = 0xC0
	_DISPLAYALLON        = 0xA5
	_DISPLAYALLON_RESUME = 0xA4
	_DISPLAYOFF          = 0xAE
	_DISPLAYON           = 0xAF
	_EXTERNALVCC         = 0x1
	_INVERTDISPLAY       = 0xA7
	_MEMORYMODE          = 0x20
	_NORMALDISPLAY       = 0xA6
	_PAGEADDR            = 0x22
	_PAGESTARTADDRESS    = 0xB0
	_SEGREMAP            = 0xA0
	_SETCOMPINS          = 0xDA
	_SETCONTRAST         = 0x81
	_SETDISPLAYCLOCKDIV  = 0xD5
	_SETDISPLAYOFFSET    = 0xD3
	_SETHIGHCOLUMN       = 0x10
	_SETLOWCOLUMN        = 0x00
	_SETMULTIPLEX        = 0xA8
	_SETPRECHARGE        = 0xD9
	_SETSEGMENTREMAP     = 0xA1
	_SETSTARTLINE        = 0x40
	_SETVCOMDETECT       = 0xDB
	_SWITCHCAPVCC        = 0x2
)

// FrameRate determines scrolling speed.
type FrameRate byte

// Possible frame rates. The value determines the number of refreshes between
// movement. The lower value, the higher speed.
const (
	FrameRate2   FrameRate = 7
	FrameRate3   FrameRate = 4
	FrameRate4   FrameRate = 5
	FrameRate5   FrameRate = 0
	FrameRate25  FrameRate = 6
	FrameRate64  FrameRate = 1
	FrameRate128 FrameRate = 2
	FrameRate256 FrameRate = 3
)

// Orientation is used for scrolling.
type Orientation byte

// Possible orientations for scrolling.
const (
	Left    Orientation = 0x27
	Right   Orientation = 0x26
	UpRight Orientation = 0x29
	UpLeft  Orientation = 0x2A
)

// DefaultOpts is the recommended default options.
var DefaultOpts = Opts{
	W:             128,
	H:             64,
	Rotated:       false,
	Sequential:    false,
	SwapTopBottom: false,
}

// Opts defines the options for the device.
type Opts struct {
	W int
	H int
	// Rotated determines if the display is rotated by 180°.
	Rotated bool
	// Sequential corresponds to the Sequential/Alternative COM pin configuration
	// in the OLED panel hardware. Try toggling this if half the rows appear to be
	// missing on your display.
	Sequential bool
	// SwapTopBottom corresponds to the Left/Right remap COM pin configuration in
	// the OLED panel hardware. Try toggling this if the top and bottom halves of
	// your display are swapped.
	SwapTopBottom bool
}

// Dev is an open handle to the display controller.
type Dev struct {
	// Communication
	c   conn.Conn
	dc  gpio.PinOut
	spi bool

	// Display size controlled by the SSD1306.
	rect image.Rectangle

	// Mutable
	// See page 25 for the GDDRAM pages structure.
	// Narrow screen will waste the end of each page.
	// Short screen will ignore the lower pages.
	// There is 8 pages, each covering an horizontal band of 8 pixels high (1
	// byte) for 128 bytes.
	// 8*128 = 1024 bytes total for 128x64 display.
	buffer []byte
	// next is lazy initialized on first Draw(). Write() skips this buffer.
	next               *image1bit.VerticalLSB
	startPage, endPage int
	startCol, endCol   int
	scrolled           bool
	halted             bool
}

func (d *Dev) String() string {
	if d.spi {
		return fmt.Sprintf("ssd1360.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
	}
	return fmt.Sprintf("ssd1360.Dev{%s, %s}", d.c, d.rect.Max)
}

// ColorModel implements display.Drawer.
//
// It is a one bit color model, as implemented by image1bit.Bit.
func (d *Dev) ColorModel() color.Model {
	return image1bit.BitModel
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements display.Drawer.
//
// It draws synchronously, once this function returns, the display is updated.
// It means that on slow bus (I²C), it may be preferable to defer Draw() calls
// to a background goroutine.
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var next []byte
	if img, ok := src.(*image1bit.VerticalLSB); ok && r == d.rect && img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		// Exact size, full frame, image1bit encoding: fast path!
		next = img.Pix
	} else {
		// Double buffering.
		if d.next == nil {
			d.next = image1bit.NewVerticalLSB(d.rect)
		}
		next = d.next.Pix
		draw.Src.Draw(d.next, r, src, sp)
	}
	return d.drawInternal(next)
}

// Write writes a buffer of pixels to the display.
//
// The format is unsual as each byte represent 8 vertical pixels at a time. The
// format is horizontal bands of 8 pixels high.
//
// This function accepts the content of image1bit.VerticalLSB.Pix.
func (d *Dev) Write(pixels []byte) (int, error) {
	if len(pixels) != len(d.buffer) {
		return 0, fmt.Errorf("ssd1306: invalid pixel stream length; expected %d bytes, got %d bytes", len(d.buffer), len(pixels))
	}
	// Write() skips d.next so it saves 1kb of RAM.
	if err := d.drawInternal(pixels); err != nil {
		return 0, err
	}
	return len(pixels), nil
}

// Scroll scrolls an horizontal band.
//
// Only one scrolling operation can happen at a time.
//
// Both startLine and endLine must be multiples of 8.
//
// Use -1 for endLine to extend to the bottom of the display.
func (d *Dev) Scroll(o Orientation, rate FrameRate, startLine, endLine int) error {
	h := d.rect.Dy()
	if endLine == -1 {
		endLine = h
	}
	if startLine >= endLine {
		return fmt.Errorf("startLine (%d) must be lower than endLine (%d)", startLine, endLine)
	}
	if startLine&7 != 0 || startLine < 0 || startLine >= h {
		return fmt.Errorf("invalid startLine %d", startLine)
	}
	if endLine&7 != 0 || endLine < 0 || endLine > h {
		return fmt.Errorf("invalid endLine %d", endLine)
	}

	startPage := uint8(startLine / 8)
	endPage := uint8(endLine / 8)
	d.scrolled = true
	if o == Left || o == Right {
		// page 28
		// <op>, dummy, <start page>, <rate>,  <end page>, <dummy>, <dummy>, <ENABLE>
		return d.sendCommand([]byte{byte(o), 0x00, startPage, byte(rate), endPage - 1, 0x00, 0xFF, 0x2F})
	}
	// page 29
	// <op>, dummy, <start page>, <rate>,  <end page>, <offset>, <ENABLE>
	// page 30: 0xA3 permits to set rows for scroll area.
	return d.sendCommand([]byte{byte(o), 0x00, startPage, byte(rate), endPage - 1, 0x01, 0x2F})
}

// StopScroll stops any scrolling previously set and resets the screen.
func (d *Dev) StopScroll() error {
	return d.sendCommand([]byte{0x2E})
}

// SetContrast changes the screen contrast.
//
// Note: values other than 0xff do not seem useful...
func (d *Dev) SetContrast(level byte) error {
	return d.sendCommand([]byte{0x81, level})
}

// SetDisplayStartLine causes the display to start from startLine, effectively
// scrolling the screen to that position.
//
// startLine must be between 0 and 63.
func (d *Dev) SetDisplayStartLine(startLine byte) error {
	if startLine > 63 {
		return fmt.Errorf("ssd1306: invalid startLine %d", startLine)
	}
	return d.sendCommand([]byte{_SETSTARTLINE | startLine})
}

// Halt turns off the display.
//
// Sending any other command afterward reenables the display.
func (d *Dev) Halt() error {
	d.halted = false
	err := d.sendCommand([]byte{0xAE})
	if err == nil {
		d.halted = true
	}
	return err
}

// Invert the display (black on white vs white on black).
func (d *Dev) Invert(blackOnWhite bool) error {
	b := []byte{0xA6}
	if blackOnWhite {
		b[0] = 0xA7
	}
	return d.sendCommand(b)
}

//

// newDev is the common initialization code that is independent of the
// communication protocol (I²C or SPI) being used.
func newDev(c conn.Conn, opts *Opts, usingSPI bool, dc gpio.PinOut) (*Dev, error) {
	if opts.W < 8 || opts.W > 128 || opts.W&7 != 0 {
		return nil, fmt.Errorf("ssd1306: invalid width %d", opts.W)
	}
	if opts.H < 8 || opts.H > 64 || opts.H&7 != 0 {
		return nil, fmt.Errorf("ssd1306: invalid height %d", opts.H)
	}

	nbPages := opts.H / 8
	pageSize := opts.W
	d := &Dev{
		c:         c,
		spi:       usingSPI,
		dc:        dc,
		rect:      image.Rect(0, 0, opts.W, opts.H),
		buffer:    make([]byte, nbPages*pageSize),
		startPage: 0,
		endPage:   nbPages,
		startCol:  0,
		endCol:    opts.W,
		// Signal that the screen must be redrawn on first draw().
		scrolled: true,
	}
	if err := d.sendCommand(getInitCmd(opts)); err != nil {
		return nil, err
	}
	return d, nil
}

func getInitCmd(opts *Opts) []byte {
	// Set COM output scan direction; C0 means normal; C8 means reversed
	comScan := byte(0xC8)
	// See page 40.
	columnAddr := byte(0xA1)
	if opts.Rotated {
		// Change order both horizontally and vertically.
		comScan = 0xC0
		columnAddr = byte(0xA0)
	}
	// See page 40.
	hwLayout := byte(0x02)
	if !opts.Sequential {
		hwLayout |= 0x10
	}
	if opts.SwapTopBottom {
		hwLayout |= 0x20
	}
	// Set the max frequency. The problem with I²C is that it creates visible
	// tear down. On SPI at high speed this is not visible. Page 23 pictures how
	// to avoid tear down. For now default to max frequency.
	freq := byte(0xF0)

	// Initialize the device by fully resetting all values.
	// Page 64 has the full recommended flow.
	// Page 28 lists all the commands.
	return []byte{
		0xAE,       // Display off
		0xD3, 0x00, // Set display offset; 0
		0x40,           // Start display start line; 0
		columnAddr,     // Set segment remap; RESET is column 127.
		comScan,        //
		0xDA, hwLayout, // Set COM pins hardware configuration; see page 40
		0x81, 0xFF, // Set max contrast
		0xA4,       // Set display to use GDDRAM content
		0xA6,       // Set normal display (0xA7 for inverted 0=lit, 1=dark)
		0xD5, freq, // Set osc frequency and divide ratio; power on reset value is 0x80.
		0x8D, 0x14, // Enable charge pump regulator; page 62
		0xD9, 0xF1, // Set pre-charge period; from adafruit driver
		0xDB, 0x40, // Set Vcomh deselect level; page 32
		0x2E,                   // Deactivate scroll
		0xA8, byte(opts.H - 1), // Set multiplex ratio (number of lines to display)
		0x20, 0x00, // Set memory addressing mode to horizontal
		0x21, 0, uint8(opts.W - 1), // Set column address (Width)
		0x22, 0, uint8(opts.H/8 - 1), // Set page address (Pages)
		0xAF, // Display on
	}
}

func (d *Dev) calculateSubset(next []byte) (int, int, int, int, bool) {
	w := d.rect.Dx()
	h := d.rect.Dy()
	startPage := 0
	endPage := h / 8
	startCol := 0
	endCol := w
	if d.scrolled {
		// Painting disable scrolling but if scrolling was enabled, this requires a
		// full screen redraw.
		d.scrolled = false
	} else {
		// Calculate the smallest square that need to be sent.
		pageSize := w

		// Top.
		for ; startPage < endPage; startPage++ {
			x := pageSize * startPage
			y := pageSize * (startPage + 1)
			if !bytes.Equal(d.buffer[x:y], next[x:y]) {
				break
			}
		}
		// Bottom.
		for ; endPage > startPage; endPage-- {
			x := pageSize * (endPage - 1)
			y := pageSize * endPage
			if !bytes.Equal(d.buffer[x:y], next[x:y]) {
				break
			}
		}
		if startPage == endPage {
			// Early exit, the image is exactly the same.
			return 0, 0, 0, 0, true
		}

		// Left.
		for ; startCol < endCol; startCol++ {
			for i := startPage; i < endPage; i++ {
				x := i*pageSize + startCol
				if d.buffer[x] != next[x] {
					goto breakLeft
				}
			}
		}
	breakLeft:

		// Right.
		for ; endCol > startCol; endCol-- {
			for i := startPage; i < endPage; i++ {
				x := i*pageSize + endCol - 1
				if d.buffer[x] != next[x] {
					goto breakRight
				}
			}
		}
	breakRight:
	}
	return startPage, endPage, startCol, endCol, false
}

// drawInternal sends image data to the controller.
func (d *Dev) drawInternal(next []byte) error {
	startPage, endPage, startCol, endCol, skip := d.calculateSubset(next)
	if skip {
		return nil
	}
	copy(d.buffer, next)

	if d.startPage != startPage || d.endPage != endPage || d.startCol != startCol || d.endCol != endCol {
		d.startPage = startPage
		d.endPage = endPage
		d.startCol = startCol
		d.endCol = endCol
	}

	pageSize := d.rect.Dx()
	for page := d.startPage; page < d.endPage; page++ {
		err := d.sendCommand([]byte{
			_PAGESTARTADDRESS | byte(page),
			_SETLOWCOLUMN | (byte(d.startCol) & 0x0F),
			_SETHIGHCOLUMN | (byte(d.startCol) >> 4),
		})
		if err != nil {
			return err
		}
		pageStart := page * pageSize
		err = d.sendData(d.buffer[pageStart+d.startCol : pageStart+d.endCol])
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dev) sendData(c []byte) error {
	if d.halted {
		// Transparently enable the display.
		if err := d.sendCommand(nil); err != nil {
			return err
		}
	}
	if d.spi {
		// 4-wire SPI.
		if err := d.dc.Out(gpio.High); err != nil {
			return err
		}
		return d.c.Tx(c, nil)
	}
	return d.c.Tx(append([]byte{i2cData}, c...), nil)
}

func (d *Dev) sendCommand(c []byte) error {
	if d.halted {
		// Transparently enable the display.
		c = append([]byte{0xAF}, c...)
		d.halted = false
	}
	if d.spi {
		if d.dc == nil {
			// 3-wire SPI.
			return errors.New("ssd1306: 3-wire SPI mode is not yet implemented")
		}
		// 4-wire SPI.
		if err := d.dc.Out(gpio.Low); err != nil {
			return err
		}
		return d.c.Tx(c, nil)
	}
	return d.c.Tx(append([]byte{i2cCmd}, c...), nil)
}

const (
	i2cCmd  = 0x00 // I²C transaction has stream of command bytes
	i2cData = 0x40 // I²C transaction has stream of data bytes
)

var _ display.Drawer = &Dev{}
