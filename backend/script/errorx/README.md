# Error Code Generation

This directory contains the configuration and tools for generating error codes in the project.

## Directory Structure

```
errorx/
├── common.yaml        # Common error codes shared across all business domains
├── metadata.yaml      # App and business domain configuration
├── evaluation.yaml    # Business-specific error codes for evaluation domain
└── code_gen.py        # Tool for generating error code files
```

## Configuration Files

### metadata.yaml
Defines the app and business domain structure:
- `app`: Application configuration
  - `name`: Application name (e.g., "cozeloop")
  - `code`: Application code (e.g., 6)
  - `business`: List of business domains
    - `name`: Business domain name
    - `code`: Business domain code

### common.yaml
Contains common error codes that are shared across all business domains:
```yaml
error_code:
  - name: CommonNoPermission
    code: 101
    message: no access permission
    no_affect_stability: true
  # ... other common error codes
```

### {biz}.yaml
Business-specific error codes (e.g., evaluation.yaml):
```yaml
error_code:
  - name: BalanceInsufficient
    code: 3001
    message: user balance is insufficient
    description: user balance is insufficient
    no_affect_stability: true
  # ... other business-specific error codes
```

## Error Code Format

Error codes are 9-digit numbers with the following structure:
```
1  2  3  4  5  6  7  8  9
_________________________
|   |  |  |  |  |  |  |  |
|app|    biz    |  sub_code |
‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
```

Where:
- `app`: Product code (cozeloop = 6)
- `biz`: Business domain code
- `sub_code`: Sub-code within the business domain

## Usage

1. Make sure you have Python 3 installed
2. Make the script executable:
   ```bash
   chmod +x code_gen.py
   ```
3. Run the script with a business domain name:
   ```bash
   # Generate code for evaluation domain
   ./code_gen.py evaluation --output-dir backend/module/evaluation/pkg/errno

   # Or use default output directory ({project_path}/backend/module/{biz}/pkg/errno)
   ./code_gen.py evaluation
   ```

The script will:
1. Validate business codes in metadata.yaml
2. Generate Go code with:
   - Constants for all error codes (common + business specific)
   - Proper error code registration in init()
   - Correct package name and imports
   - Proper formatting and comments

## Adding New Business Domains

1. Add the business domain to metadata.yaml:
   ```yaml
   business:
     - name: your_biz
       code: 123  # Must be unique
   ```

2. Create a new {biz}.yaml file with business-specific error codes:
   ```yaml
   error_code:
     - name: YourError
       code: 1001
       message: error message
       description: error description
       no_affect_stability: true
   ```

3. Run the generator:
   ```bash
   ./code_gen.py your_biz
   ```

## Adding New Error Codes

1. For common errors, add to common.yaml
2. For business-specific errors, add to the corresponding {biz}.yaml
3. Run the generator to update the code
