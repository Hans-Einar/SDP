"""SDUI 0.2 frontend and structural console preview; no SDL/GUI runtime."""
from .ast import SduiError, to_data
from .parser import parse
from .validate import validate

__all__ = ['SduiError', 'to_data', 'parse', 'validate']
