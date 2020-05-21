package fbsparser

import (
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/test/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestFbsSchemaParser(t *testing.T) {
	schema, err := ParseSchemaFile("non-existent.fbs")
	assert.True(t, schema == nil)
	assert.Err(t, err)

	file, err := ioutil.TempFile("", "fbs-test")
	assert.NoErr(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(`// Example schema file
// Note the property comments - they test multiple different formats

namespace FbsC;

enum Planet:byte { Mercury = 0, Venus, Earth = 2 }

table Being {
  age:short = 150;
  health:short = 100; //< Current health points
  name:string; // Full being name
  friendly:bool = false (deprecated);
  location:Planet = Earth;

  /// All worldly belongings of this being
  belongings:[Item];
}

table Item {
  name:string;
  weight:short;
}

root_type Being;`)
	assert.NoErr(t, err)

	schema, err = ParseSchemaFile(file.Name())
	assert.NoErr(t, err)
	assert.True(t, schema != nil)

	assert.Eq(t, 1, schema.EnumsLength())
	assert.Eq(t, 2, schema.ObjectsLength())

	var enum *reflection.Enum
	assert.True(t, schema.Enums(enum, 0))
	assert.True(t, enum != nil)
	assert.Eq(t, "Planet", string(enum.Name()))

}
