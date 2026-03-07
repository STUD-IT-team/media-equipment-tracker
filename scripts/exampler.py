#!/usr/bin/env python3
import os
import re
from pathlib import Path
import argparse
from typing import Tuple, Optional

class BaseKeyValueExtensionParser:
    def parse(self, line: str) -> Tuple[str, str]:
        """Parse a line into key-value pair.
        
        Args:
            line: The line to parse
            
        Returns:
            Tuple of (key, value). For non-key-value lines, returns (line, '')
        """
        raise NotImplementedError


class EnvKeyValueExtensionParser(BaseKeyValueExtensionParser):
    def parse(self, line: str) -> Tuple[str, str]:
        line = line.strip()
        if '=' in line:
            key, value = line.split('=', 1)
            return key.strip(), value.strip()
        return line, ''


class YamlKeyValueExtensionParser(BaseKeyValueExtensionParser):
    def parse(self, line: str) -> Tuple[str, str]:
        if ':' in line and not line.strip().startswith('#'):
            parts = line.split(':', 1)
            return parts[0].strip(), parts[1].strip()
        return line, ''


class BaseKeyValueWriter:
    def write(self, key: str, value: str) -> str:
        """Format a key-value pair into a line.
        
        Args:
            key: The key
            value: The value
            
        Returns:
            Formatted line
        """
        raise NotImplementedError


class EnvKeyValueWriter(BaseKeyValueWriter):
    def write(self, key: str, value: str) -> str:
        return f"{key}={value}"


class YamlKeyValueWriter(BaseKeyValueWriter):
    def write(self, key: str, value: str) -> str:
        return f"{key}: {value}"


class FileExtensionSolution:
    """Factory class for getting appropriate parsers and writers."""
    
    _EXTENSION_MAP = {
        '.env': (EnvKeyValueExtensionParser, EnvKeyValueWriter),
        '.yaml': (YamlKeyValueExtensionParser, YamlKeyValueWriter),
        '.yml': (YamlKeyValueExtensionParser, YamlKeyValueWriter),
    }

    @classmethod
    def get_parser(cls, file_extension: str) -> BaseKeyValueExtensionParser:
        """Get parser for given file extension.
        
        Args:
            file_extension: The file extension (e.g. '.env')
            
        Returns:
            Appropriate parser instance
            
        Raises:
            ValueError: If extension is not supported
        """
        if file_extension not in cls._EXTENSION_MAP:
            raise ValueError(f"Unsupported file extension: {file_extension}")
        return cls._EXTENSION_MAP[file_extension][0]()

    @classmethod
    def get_writer(cls, file_extension: str) -> BaseKeyValueWriter:
        """Get writer for given file extension.
        
        Args:
            file_extension: The file extension (e.g. '.env')
            
        Returns:
            Appropriate writer instance
            
        Raises:
            ValueError: If extension is not supported
        """
        if file_extension not in cls._EXTENSION_MAP:
            raise ValueError(f"Unsupported file extension: {file_extension}")
        return cls._EXTENSION_MAP[file_extension][1]()

    @classmethod
    def supported_extensions(cls) -> list:
        """Get list of supported file extensions."""
        return list(cls._EXTENSION_MAP.keys())


def sanitize_value(key: str, value: str) -> str:
    """Generate example value based on key name and original value.
    
    Args:
        key: The configuration key
        value: The original value (for type detection)
        
    Returns:
        Example value string
    """
    key_lower = key.lower()
    
    # Common patterns based on key name
    if any(x in key_lower for x in ['secret', 'password', 'token', 'key']):
        return f'your_{key_lower}_here'
    if 'email' in key_lower or 'mail' in key_lower:
        return 'your_email@example.com'
    if 'url' in key_lower or 'uri' in key_lower:
        return 'https://example.com/path'
    if 'port' in key_lower:
        return '8080'
    if 'host' in key_lower or 'server' in key_lower:
        return 'example.com'
    if 'database' in key_lower or 'db' in key_lower:
        return 'your_database_name'
    
    # Type detection from original value
    if value.isdigit():
        return '12345'
    if value.replace('.', '').isdigit():
        return '12.34'
    if value.lower() in ('true', 'false'):
        return 'true_or_false'
    
    return 'example_value'


def get_output_path(file_path: str, suffix: str = '.example') -> str:
    """Generate output path with pattern: filename.example.extension.
    
    Args:
        file_path: Original file path
        suffix: Suffix to insert before extension
        
    Returns:
        New file path with suffix inserted
    """
    path = Path(file_path)
    if path.name.startswith('.'):  # Handle dotfiles
        return str(path.parent / f"{path.name}{suffix}{path.suffix}")
    return str(path.parent / f"{path.stem}{suffix}{path.suffix}")


def process_file(
    file_path: str,
    suffix: str = '.example',
    parser: Optional[BaseKeyValueExtensionParser] = None,
    writer: Optional[BaseKeyValueWriter] = None,
    override: bool = False
) -> bool:
    """Process a file to create its example version.
    
    Args:
        file_path: Path to original file
        suffix: Suffix for example file
        parser: Parser instance
        writer: Writer instance
        override: Whether to overwrite existing example files
        
    Returns:
        True if file was processed, False if skipped
        
    Raises:
        ValueError: If parser or writer is not provided
    """
    if parser is None or writer is None:
        raise ValueError("Both parser and writer must be provided")
    
    if '.example.' in file_path.lower():
        print(f"[SKIP] Already an example file: {file_path}")
        return False
    
    output_path = get_output_path(file_path, suffix)
    
    if os.path.exists(output_path) and not override:
        print(f"[SKIP] Example file exists (use --override to overwrite): {output_path}")
        return False
    
    print(f"[PROCESS] {file_path} -> {output_path}")
    
    try:
        with open(file_path, 'r') as f:
            content = f.read()
        
        new_lines = []
        for line in content.split('\n'):
            stripped_left = line.lstrip()
            spaces = len(line) - len(stripped_left)
            stripped = stripped_left.strip()
            
            # Preserve comments and empty lines
            if not stripped or stripped.startswith('#'):
                new_lines.append(line)
                continue
            
            key, value = parser.parse(line)
            if value:
                sanitized = sanitize_value(key, value)
                new_lines.append(" " * spaces + writer.write(key, sanitized))
            else:
                new_lines.append(line)
        
        with open(output_path, 'w') as f:
            f.write('\n'.join(new_lines))
        return True
            
    except Exception as e:
        print(f"[ERROR] Processing {file_path}: {str(e)}")
        return False


def should_exclude(filepath: str, exclude_patterns: list) -> bool:
    """Check if file should be excluded based on patterns.
    
    Args:
        filepath: Path to check
        exclude_patterns: List of patterns to exclude
        
    Returns:
        True if file should be excluded
    """
    filepath_lower = filepath.lower()
    return any(pattern.lower() in filepath_lower for pattern in exclude_patterns)


def find_config_files(
    directories: list,
    extensions: list,
    exclude_patterns: list
) -> list:
    """Find all config files in specified directories.
    
    Args:
        directories: Directories to search
        extensions: File extensions to include
        exclude_patterns: Patterns to exclude
        
    Returns:
        List of matching file paths
    """
    config_files = []
    for directory in directories:
        for root, _, files in os.walk(directory):
            for file in files:
                filepath = os.path.join(root, file)
                if (any(file.lower().endswith(ext) for ext in extensions) and 
                    not should_exclude(filepath, exclude_patterns)):
                    config_files.append(filepath)
    return config_files


def main():
    parser = argparse.ArgumentParser(
        description='Generate example config files from real ones.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''Examples:
  python3 generate_examples.py --dirs /path/to/configs
  python3 generate_examples.py --dirs /path1 /path2 --exclude test temp
  python3 generate_examples.py --dirs . --extensions .env .yaml --suffix .template --override'''
    )
    parser.add_argument('--dirs', nargs='+', required=True,
                       help='Directories to search for config files')
    parser.add_argument('--extensions', nargs='+',
                       default=FileExtensionSolution.supported_extensions(),
                       help='File extensions to consider as config files')
    parser.add_argument('--suffix', default='.example',
                       help='Suffix to add before extension (default: .example)')
    parser.add_argument('--exclude', nargs='+', default=[],
                       help='Case-insensitive patterns to exclude from processing')
    parser.add_argument('--override', action='store_true', default=False,
                       help='Overwrite existing example files')
    
    args = parser.parse_args()
    
    # Verify directories exist
    for directory in args.dirs:
        if not os.path.isdir(directory):
            print(f"[ERROR] Directory not found: {directory}")
            return
    
    # Find all config files
    config_files = find_config_files(args.dirs, args.extensions, args.exclude)
    
    if not config_files:
        print("[INFO] No config files found with the specified extensions.")
        return
    
    print(f"[INFO] Found {len(config_files)} config files to process")
    
    # Process each file
    processed_count = 0
    skipped_count = 0
    error_count = 0
    
    for file in config_files:
        file_extension = Path(file).suffix.lower()
        try:
            parser = FileExtensionSolution.get_parser(file_extension)
            writer = FileExtensionSolution.get_writer(file_extension)
            if process_file(file, args.suffix, parser, writer, args.override):
                processed_count += 1
            else:
                skipped_count += 1
        except ValueError as e:
            print(f"[SKIP] {file}: {str(e)}")
            skipped_count += 1
        except Exception as e:
            print(f"[ERROR] Unexpected error with {file}: {str(e)}")
            error_count += 1
    
    print("\n[SUMMARY]")
    print(f"Processed files: {processed_count}")
    print(f"Skipped files: {skipped_count}")
    if error_count > 0:
        print(f"Errors encountered: {error_count}")


if __name__ == '__main__':
    main()