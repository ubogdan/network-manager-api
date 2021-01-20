package license

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextValidPeriod(t *testing.T) {
	tt := []struct {
		Created  time.Time
		Expire   time.Time
		Now      time.Time
		Result   time.Time
		Validity time.Duration
		isError  bool
	}{
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1547078400, 0), // 10 Jan 2019
			Expire:   time.Unix(1577836800, 0), // 1 Jan 2020
			Result:   time.Unix(1546300800, 0), // 1 jan 2019
			Validity: 30 * 24 * time.Hour,      // 30 days
		},
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1548806400, 0), // 30 Jan 2019
			Expire:   time.Unix(1577836800, 0), // 1 Jan 2020
			Result:   time.Unix(1546300800, 0), // 1 jan 2019
			Validity: 30 * 24 * time.Hour,      // 30 days
		},
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1550188800, 0), // 15 Feb 2019
			Expire:   time.Unix(1577836800, 0), // 1 Jan 2020
			Result:   time.Unix(1548892800, 0), // 31 Jan 2019
			Validity: 30 * 24 * time.Hour,      // 30 days
		},
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1551571200, 0), // 5 March 2019
			Expire:   time.Unix(1577836800, 0), // 1 Jan 2020
			Result:   time.Unix(1551484800, 0), // 2 March 2019
			Validity: 30 * 24 * time.Hour,      // 30 days
		},
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1575504000, 0), // 5 Dec 2019
			Expire:   time.Unix(1577836800, 0), // 1 Jan 2020
			Result:   time.Unix(1574812800, 0), // 27 Nov 2019
			Validity: 30 * 24 * time.Hour,      // 30 days
		},
		{
			Created:  time.Unix(1546300800, 0), // 1 jan 2019
			Now:      time.Unix(1577664000, 0), // 30 Dec 2019 License expired
			Expire:   time.Unix(1577836800, 0), // 27 Jan 2020
			isError:  true,
			Validity: 30 * 24 * time.Hour, // 30 days
		},
	}

	for _, test := range tt {
		period, err := nextValidPeriod(test.Created, test.Expire, test.Now, test.Validity)
		if test.isError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, period, test.Result)
	}
}
