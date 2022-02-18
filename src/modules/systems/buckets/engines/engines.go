package engines

import (
	_ "herbbuckets/modules/systems/buckets/engines/localbucket" //local bucket driver
	_ "herbbuckets/modules/systems/buckets/engines/s3bucket"    //s3 bucket driver
)
