import { frequencyToSpeed } from "./freq-to-speed.function";

describe("frequencyToSpeed", () => {
  it("should interpolate 0 Hz", () => {
    // GIVEN
    const freq = 0.0;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(0.0, 0.0)
  });

  it("should interpolate 0.75 Hz", () => {
    // GIVEN
    const freq = 0.75;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(1.35, 0.01)
  });

  it("should interpolate 2.0 Hz", () => {
    // GIVEN
    const freq = 2.0;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(2.0, 0.01)
  });

  it("should interpolate 2.5 Hz", () => {
    // GIVEN
    const freq = 2.5;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(1.25, 0.01)
  });


  it("should interpolate 3.0 Hz", () => {
    // GIVEN
    const freq = 3.0;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(0.0, 0.01)
  });
  it("should interpolate high frequency", () => {
    // GIVEN
    const freq = 100000.0;
    // WHEN
    const speed = frequencyToSpeed(freq);
    // THEN
    expect(speed).toBeCloseTo(0.0, 0.01)
  });
});
