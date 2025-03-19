package redis

import (
	"go-fiber-api/internal/core/response"
	"testing"

	"go-fiber-api/internal/core/config"

	"github.com/stretchr/testify/assert"
)

func TestRedisClient_Set_Byte(t *testing.T) {
	cfg := &config.Configuration{
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDB:       10,
	}

	tests := []struct {
		name          string
		key           string
		value         []byte
		ttl           []int
		wantErr       bool
		expectedError string
	}{
		{
			name:          "success",
			key:           "test-key-success-case-byte",
			value:         []byte("test-value"),
			ttl:           []int{5},
			wantErr:       false,
			expectedError: "",
		},
		{
			name:          "empty key",
			key:           "",
			value:         []byte("c"),
			ttl:           []int{5},
			wantErr:       true,
			expectedError: "key cannot be empty",
		},
		{
			name:          "empty value",
			key:           "test-key-empty-value-byte",
			value:         []byte(""),
			ttl:           []int{5},
			wantErr:       false,
			expectedError: "",
		},
		{
			name:          "negative ttl",
			key:           "test-key-negative-ttl-byte",
			value:         []byte("test-value"),
			ttl:           []int{-1},
			wantErr:       true,
			expectedError: "TTL must be greater than 0",
		},
		{
			name:          "no ttl",
			key:           "test-key-no-ttl-byte",
			value:         []byte("test-value"),
			ttl:           []int{},
			wantErr:       false,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := ProvideClient(cfg)
			c := client.(*clientImpl)
			defer ResetProvideClient()
			err := c.Set(tt.key, tt.value, tt.ttl...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRedisClient_Set_String(t *testing.T) {
	cfg := &config.Configuration{
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDB:       10,
	}

	tests := []struct {
		name          string
		key           string
		value         any
		ttl           []int
		wantErr       bool
		expectedError string
	}{
		{
			name:          "success",
			key:           "test-key-success-case",
			value:         interface{}("test-value"),
			ttl:           []int{5},
			wantErr:       false,
			expectedError: "",
		},
		{
			name:          "marshal error",
			key:           "test-key-will-fail",
			value:         make(chan int), // invalid value for JSON marshaling
			ttl:           []int{60},
			wantErr:       true,
			expectedError: "failed to marshal value: json: unsupported type: chan int",
		},
		{
			name:          "empty key",
			key:           "",
			value:         "test-value",
			ttl:           []int{5},
			wantErr:       true,
			expectedError: "key cannot be empty",
		},
		{
			name:          "empty value",
			key:           "test-key-empty-value",
			value:         "",
			ttl:           []int{5},
			wantErr:       false,
			expectedError: "",
		},
		{
			name:          "negative ttl",
			key:           "test-key-negative-ttl",
			value:         "test-value",
			ttl:           []int{-1},
			wantErr:       true,
			expectedError: "TTL must be greater than 0",
		},
		{
			name:          "no ttl",
			key:           "test-key-no-ttl",
			value:         "test-value",
			ttl:           []int{},
			wantErr:       false,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := ProvideClient(cfg)
			c := client.(*clientImpl)
			defer ResetProvideClient()

			err := c.Set(tt.key, tt.value, tt.ttl...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)

				var val []byte
				err := c.Get(tt.key, &val)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, val)
			}
		})
	}
}

func TestRedisClient_Set_Model(t *testing.T) {
	cfg := &config.Configuration{
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDB:       10,
	}

	tests := []struct {
		name  string
		key   string
		value any

		ttl           []int
		wantErr       bool
		expectedError string
	}{

		{
			name:          "cast model",
			key:           "test-key-cast-model",
			value:         response.ResponseDTO{Message: "test-message"},
			ttl:           []int{10},
			wantErr:       false,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := ProvideClient(cfg)
			c := client.(*clientImpl)
			defer ResetProvideClient()

			err := c.Set(tt.key, tt.value, tt.ttl...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)

				var val response.ResponseDTO
				err := c.Get(tt.key, &val)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, val)
			}
		})
	}
}

func TestRedisClient_Del(t *testing.T) {
	cfg := &config.Configuration{
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDB:       10,
	}

	tests := []struct {
		name          string
		key           string
		wantErr       bool
		expectedError string
	}{
		{
			name:    "when_del_with_key_success_should_got_nil_error",
			key:     "test-key-success-case",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := ProvideClient(cfg)
			c := client.(*clientImpl)
			defer ResetProvideClient()

			err := c.Del(tt.key)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestRedisClient_Close(t *testing.T) {
	cfg := &config.Configuration{
		RedisHost:     "localhost",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDB:       10,
	}

	tests := []struct {
		name          string
		key           string
		wantErr       bool
		expectedError string
	}{
		{
			name:    "when_close_client_success_should_got_nil_error",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := ProvideClient(cfg)
			c := client.(*clientImpl)
			defer ResetProvideClient()

			err := c.Close()

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
