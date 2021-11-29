package rateLimit

import "time"

type nowFactory func() time.Time
