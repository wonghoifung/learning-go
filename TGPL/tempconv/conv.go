package tempconv
import "fmt"
func init() {
  fmt.Println("conv init")
  fmt.Println(BoilingC.String())
}
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c * 9 / 5 + 32) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }
