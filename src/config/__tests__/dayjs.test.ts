import { describe, it, expect } from 'vitest';
import dayjs from '../dayjs';

describe('dayjs', () => {
  it('deve estar configurado com locale pt-br', () => {
    expect(dayjs.locale()).toBe('pt-br');
  });

  it('deve ter plugin customParseFormat disponível', () => {
    const date = dayjs('2024-01-01', 'YYYY-MM-DD');
    expect(date.isValid()).toBe(true);
  });

  it('deve ter plugin weekOfYear disponível', () => {
    const date = dayjs('2024-01-01');
    expect(date.week()).toBeDefined();
    expect(typeof date.week()).toBe('number');
  });

  it('deve ter plugin advancedFormat disponível', () => {
    const date = dayjs('2024-01-01');
    expect(date.format('Q')).toBeDefined();
  });

  it('deve ter plugin isSameOrAfter disponível', () => {
    const date1 = dayjs('2024-01-01');
    const date2 = dayjs('2024-01-02');
    expect(date2.isSameOrAfter(date1)).toBe(true);
  });

  it('deve ter plugin isSameOrBefore disponível', () => {
    const date1 = dayjs('2024-01-01');
    const date2 = dayjs('2024-01-02');
    expect(date1.isSameOrBefore(date2)).toBe(true);
  });

  it('deve ter plugin isBetween disponível', () => {
    const date = dayjs('2024-01-15');
    const start = dayjs('2024-01-01');
    const end = dayjs('2024-01-31');
    expect(date.isBetween(start, end)).toBe(true);
  });

  it('deve formatar datas corretamente', () => {
    const date = dayjs('2024-01-01');
    expect(date.format('YYYY-MM-DD')).toBe('2024-01-01');
  });
});




