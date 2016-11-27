package missilelauncher

/*
#cgo CFLAGS: -std=gnu99
#cgo linux  pkg-config: libmissilelauncher
#include <stdio.h>
#include <stdlib.h>
#include <libmissilelauncher/libmissilelauncher.h>
*/
import "C"
import "unsafe"
import "runtime"
import "errors"

/*
 * Types
 */
type Launcher C.ml_launcher_t

type LauncherArray struct {
	array **C.ml_launcher_t
	size  C.uint32_t
}

type LauncherType C.ml_launcher_type

const (
	ML_NOT_LAUNCHER      LauncherType = C.ML_NOT_LAUNCHER
	ML_STANDARD_LAUNCHER              = C.ML_STANDARD_LAUNCHER
)

type LauncherDirection C.ml_launcher_direction

const (
	ML_DOWN  LauncherDirection = C.ML_DOWN
	ML_UP                      = C.ML_UP
	ML_LEFT                    = C.ML_LEFT
	ML_RIGHT                   = C.ML_RIGHT
)

/*
 * Utilities
 */
func ecToError(ec C.ml_error_code) error {
	if ec == C.ML_OK {
		return nil
	}

	str := C.GoString(C.ml_error_to_str(ec))
	return errors.New(str)
}

/*
 * Library
 */
func LibraryInit() error {
	return ecToError(C.ml_library_init())
}

func LibraryCleanup() error {
	return ecToError(C.ml_library_cleanup())
}

func LibaryIsInit() bool {
	if C.ml_library_is_init() == 0 {
		return false
	} else {
		return true
	}
}

/*
 * Launcher Array
 */
func NewLauncherArray() (error, *LauncherArray) {
	la := &LauncherArray{
		array: nil,
		size:  0,
	}

	runtime.SetFinalizer(la, freeLauncherArray)

	err := ecToError(C.ml_launcher_array_new(&la.array, &la.size))
	if err != nil {
		return err, nil
	}

	return nil, la
}

func freeLauncherArray(la *LauncherArray) {
	if la.array != nil {
		C.ml_launcher_array_free(la.array)
	}
}

func (la *LauncherArray) ToSlice() []*Launcher {
	return (*[1 << 30]*Launcher)(unsafe.Pointer(la.array))[:la.size:la.size]
}

/*
 * Launcher
 */
func (l *Launcher) Ref() error {
	return ecToError(C.ml_launcher_reference(l))
}

func (l *Launcher) Deref() error {
	return ecToError(C.ml_launcher_dereference(l))
}

func (l *Launcher) Claim() error {
	return ecToError(C.ml_launcher_claim(l))
}

func (l *Launcher) Unclaim() error {
	return ecToError(C.ml_launcher_unclaim(l))
}

func (l *Launcher) Type() LauncherType {
	return LauncherType(C.ml_launcher_get_type(l))
}

func (l *Launcher) Fire() error {
	return ecToError(C.ml_launcher_fire(l))
}

func (l *Launcher) Move(ld LauncherDirection) error {
	return ecToError(C.ml_launcher_move(l,
		C.ml_launcher_direction(ld)))
}

func (l *Launcher) Stop() error {
	return ecToError(C.ml_launcher_stop(l))
}

func (l *Launcher) MoveMSeconds(ld LauncherDirection, msec uint32) error {
	return ecToError(C.ml_launcher_move_mseconds(l,
		C.ml_launcher_direction(ld), C.uint32_t(msec)))
}

func (l *Launcher) Zero() error {
	return ecToError(C.ml_launcher_zero(l))
}

func (l *Launcher) LEDOn() error {
	return ecToError(C.ml_launcher_led_on(l))
}

func (l *Launcher) LEDOff() error {
	return ecToError(C.ml_launcher_led_off(l))
}

func (l *Launcher) GetLEDState() bool {
	if C.ml_launcher_get_led_state(l) == 0 {
		return false
	} else {
		return true
	}
}
