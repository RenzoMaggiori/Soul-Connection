export function getDarkerColor(hexColor: string, percent: number = 20): string {
    const hex = hexColor.replace('#', '');
    const rgb = hex.match(/.{2}/g)?.map(val => parseInt(val, 16)) ?? [0, 0, 0];
    const darkerRgb = rgb.map(val => Math.floor(val * (100 - percent) / 100));
    return `#${darkerRgb.map(val => val.toString(16).padStart(2, '0')).join('')}`;
}