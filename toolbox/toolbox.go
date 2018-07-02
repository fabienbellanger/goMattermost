package toolbox

import (
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// CheckError : Gestion des erreurs
func CheckError(err error, exitCode int) {
	if err != nil {
		if exitCode != 0 {
			color.Red("Error (" + strconv.Itoa(exitCode) + "): " + err.Error() + "\n\n")

			os.Exit(exitCode)
		} else {
			panic(err.Error())
		}
	}
}

// Ucfirst : Met la première lettre d'une chaîne de caractères en majuscule
func Ucfirst(s string) string {
	sToUnicode := []rune(s) // Tableau de caractères Unicode pour gérér les caractères accentués

	return strings.ToUpper(string(sToUnicode[0])) + string(sToUnicode[1:])
}
