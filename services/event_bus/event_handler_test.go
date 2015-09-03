package event_bus

import (
	"testing"
	"fmt"
	"math/rand"
	"io/ioutil"

	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"

	"github.com/QubitProducts/bamboo/configuration"
)

// Yanked off http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6
    letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
    b := make([]byte, n)
    // A rand.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = rand.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}

func orPanic(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error! %s", err))
	}
}

func TestIsReloadRequired(t *testing.T) {
	var config configuration.Configuration
	_ = config

	Convey("#isReloadRequired", t, func () {

		tmpPath := fmt.Sprintf("/tmp/bamboo_irr%v.conf", rand.Int31())
		contents := randString(int(rand.Int31n(1048576)))

		orPanic(ioutil.WriteFile(tmpPath, []byte(contents), 0644))

		Convey("When we test with a file's own contents", func() {

			required, err := isReloadRequired(tmpPath, contents)
			orPanic(err)

			Convey("a reload should not be required", func() {
				So(required, ShouldEqual, false)
			})
		})

		Convey("When we test with changed config", func() {
			required, err := isReloadRequired(tmpPath, contents + "arst")
			orPanic(err)

			Convey("a reload should be required", func() {
				So(required, ShouldEqual, true)
			})
		})

		Convey("When we test against an absent file", func() {
			required, err := isReloadRequired("/tmp/bamboo_irr_nonexistant.conf", contents)
			orPanic(err)

			Convey("a reload should be required", func() {
				So(required, ShouldEqual, true)
			})
		})

		Convey("When we test against an invalid file", func() {
			_, err := isReloadRequired("/dev/null/foobar", contents)

			Convey("we should get an error", func() {
				So(err, ShouldNotEqual, nil)
			})
		})
	})
}
