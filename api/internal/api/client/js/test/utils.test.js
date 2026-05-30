import { describe, it, expect } from 'vitest';
import { ohMySimpleHash } from '../utils.js';

describe('ohMySimpleHash', () => {
  it('returns 0 for undefined', () => {
    expect(ohMySimpleHash(undefined)).toBe(0);
  });

  it('returns 0 for empty string', () => {
    expect(ohMySimpleHash('')).toBe(0);
  });

  it('is deterministic for same input', () => {
    expect(ohMySimpleHash('hello world')).toBe(ohMySimpleHash('hello world'));
  });

  it('differs for different inputs', () => {
    expect(ohMySimpleHash('foo')).not.toBe(ohMySimpleHash('bar'));
  });

  it('returns a number', () => {
    expect(typeof ohMySimpleHash('test')).toBe('number');
  });

  it('handles long strings', () => {
    const long = 'x'.repeat(10000);
    expect(typeof ohMySimpleHash(long)).toBe('number');
    expect(ohMySimpleHash(long)).toBe(ohMySimpleHash(long));
  });

  it('is sensitive to character order', () => {
    expect(ohMySimpleHash('ab')).not.toBe(ohMySimpleHash('ba'));
  });
});
