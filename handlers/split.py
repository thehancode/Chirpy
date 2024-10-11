
import re
import os
import sys
package_declaration = """package handlers

import (
\t"Chirpy/handlers/models"
\t"Chirpy/internal/auth"
\t"Chirpy/internal/database"
\t"database/sql"
\t"encoding/json"
\t"fmt"
\t"log"
\t"net/http"
\t"time"

\t"github.com/google/uuid"
\t"github.com/gorilla/mux"
)

"""
def split_handlers(input_file, output_dir):
    # Ensure the output directory exists
    os.makedirs(output_dir, exist_ok=True)
    
    # Regular expression to match the handler function definition
    # Example: func (cfg *ApiConfig) PostUserHandler(w http.ResponseWriter, r *http.Request) {
    func_regex = re.compile(r'^func\s+\(cfg\s+\*\w+\)\s+(\w+)\s*\(w\s+http\.ResponseWriter,\s+r\s+\*http\.Request\)\s*\{')
    
    with open(input_file, 'r') as f:
        lines = f.readlines()
    
    inside_function = False
    brace_count = 0
    current_func_name = ""
    current_func_lines = []
    
    for idx, line in enumerate(lines):
        if not inside_function:
            match = func_regex.match(line.strip())
            if match:
                inside_function = True
                current_func_name = match.group(1)
                current_func_lines = [line]
                # Count opening braces in the current line
                brace_count = line.count('{') - line.count('}')
                continue
        else:
            current_func_lines.append(line)
            # Update brace count
            brace_count += line.count('{') - line.count('}')
            if brace_count == 0:
                inside_function = False
                func_filename = f"{current_func_name}.go"
                func_filepath = os.path.join(output_dir, func_filename)
                with open(func_filepath, 'w') as func_file:
                    func_file.write(package_declaration)  
                    func_file.writelines(current_func_lines)
                print(f"Extracted {current_func_name} to {func_filename}")
                current_func_name = ""
                current_func_lines = []
    
    if inside_function:
        print("Warning: File ended while still inside a function.")
                

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python split_handlers.py <input_file.go> <output_directory>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_dir = sys.argv[2]
    
    split_handlers(input_file, output_dir)
