package envs_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/nyaruka/gocommon/jsonx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentMarshaling(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	// can't create with invalid date format
	_, err = envs.ReadEnvironment(json.RawMessage(`{"date_format": "YYYYYYYYYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid time format
	_, err = envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid language
	_, err = envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "allowed_languages": ["elvish"]}`))
	assert.Error(t, err)

	// can't create with invalid country
	_, err = envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "default_country": "Narnia"}`))
	assert.Error(t, err)

	// can't create with invalid timzeone
	_, err = envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Cuenca"}`))
	assert.Error(t, err)

	// empty environment uses all defaults
	env, err := envs.ReadEnvironment(json.RawMessage(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, envs.DateFormatYearMonthDay, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinute, env.TimeFormat())
	assert.Equal(t, envs.DefaultNumberFormat, env.NumberFormat())
	assert.Equal(t, envs.NilLanguage, env.DefaultLanguage())
	assert.Nil(t, env.AllowedLanguages())
	assert.Equal(t, envs.NilCountry, env.DefaultCountry())
	assert.Equal(t, 640, env.MaxValueLength())
	assert.Nil(t, env.LocationResolver())

	// can create with valid values
	env, err = envs.ReadEnvironment(json.RawMessage(`{
		"date_format": "DD-MM-YYYY", 
		"time_format": "tt:mm:ss", 
		"allowed_languages": ["eng", "fra"], 
		"default_country": "RW", 
		"timezone": "Africa/Kigali"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, envs.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, envs.Language("eng"), env.DefaultLanguage())
	assert.Equal(t, []envs.Language{envs.Language("eng"), envs.Language("fra")}, env.AllowedLanguages())
	assert.Equal(t, envs.Country("RW"), env.DefaultCountry())
	assert.Equal(t, "en-RW", env.DefaultLocale().ToBCP47())
	assert.Nil(t, env.LocationResolver())

	data, err := jsonx.Marshal(env)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"date_format":"DD-MM-YYYY","time_format":"tt:mm:ss","timezone":"Africa/Kigali","allowed_languages":["eng","fra"],"number_format":{"decimal_symbol":".","digit_grouping_symbol":","},"default_country":"RW","redaction_policy":"none","max_value_length":640}`)
}

func TestEnvironmentEqual(t *testing.T) {
	env1, err := envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	require.NoError(t, err)

	env2, err := envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	require.NoError(t, err)

	env3, err := envs.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kampala"}`))
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
	assert.True(t, env2.Equal(env1))
	assert.False(t, env1.Equal(env3))

	// marshal and unmarshal env 1 again
	env1JSON, err := jsonx.Marshal(env1)
	require.NoError(t, err)
	env1, err = envs.ReadEnvironment(env1JSON)
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
}

func TestEnvironmentBuilder(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	env := envs.NewBuilder().
		WithDateFormat(envs.DateFormatDayMonthYear).
		WithTimeFormat(envs.TimeFormatHourMinuteSecond).
		WithTimezone(kgl).
		WithAllowedLanguages([]envs.Language{envs.Language("fra"), envs.Language("eng")}).
		WithDefaultCountry(envs.Country("RW")).
		WithNumberFormat(&envs.NumberFormat{DecimalSymbol: "'"}).
		WithRedactionPolicy(envs.RedactionPolicyURNs).
		WithMaxValueLength(1024).
		Build()

	assert.Equal(t, envs.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, []envs.Language{envs.Language("fra"), envs.Language("eng")}, env.AllowedLanguages())
	assert.Equal(t, envs.Country("RW"), env.DefaultCountry())
	assert.Equal(t, &envs.NumberFormat{DecimalSymbol: "'"}, env.NumberFormat())
	assert.Equal(t, envs.RedactionPolicyURNs, env.RedactionPolicy())
	assert.Equal(t, 1024, env.MaxValueLength())
	assert.Nil(t, env.LocationResolver())
}
