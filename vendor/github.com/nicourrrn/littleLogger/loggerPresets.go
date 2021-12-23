package littleLogger

import (
	"fmt"
	"time"
)

func FormatterClassic() string{
	return fmt.Sprintf("%s: $msg", time.Now().String())
}

func FormatterMinimal() string{
	return "$msg"
}
