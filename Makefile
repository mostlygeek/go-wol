# Define the target architecture and build directory
BUILD_DIR := build

# Default target to build all binaries
all: $(BUILD_DIR)/wol_http_linux_amd64

# Ensure the build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Example binary build targets
# Replace these with your actual binary build commands
$(BUILD_DIR)/wol_http_linux_amd64: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/wol_http_linux_amd64 main.go

# Add more binary build targets as needed
