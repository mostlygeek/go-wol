# Define the target architecture and build directory
TARGET_ARCH := linux-amd64
BUILD_DIR := build

# Ensure the build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Example binary build targets
# Replace these with your actual binary build commands
$(BUILD_DIR)/wol_http_$(GOOS)_$(GOARCH): $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/wol_http_linux_amd64 ./path/to/source

# Add more binary build targets as needed
