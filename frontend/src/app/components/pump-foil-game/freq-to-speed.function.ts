export function frequencyToSpeed(freq: number /* Hz*/): number {
  // Find index of entry in table where freq is above the input freq
  const index = HZ_TO_SPEED_TRANSLATION_TABLE.findIndex(([tableFreq, speed]) => tableFreq > freq);
  const prevIndex = Math.max(0, index - 1);
  console.log("index:", index, "prevIndex:", prevIndex);
  const speedLow = HZ_TO_SPEED_TRANSLATION_TABLE[prevIndex][1];
  const speedHigh = HZ_TO_SPEED_TRANSLATION_TABLE[index][1];

  // Interpolate speed between indices index and prevIndex
  if (index <= 0 || index >= HZ_TO_SPEED_TRANSLATION_TABLE.length) {
    // If the index is out of bounds, return the lowest or highest speed in the table
    return HZ_TO_SPEED_TRANSLATION_TABLE[Math.max(0, Math.min(index, HZ_TO_SPEED_TRANSLATION_TABLE.length - 1))][1];
  }

  const freqLow = HZ_TO_SPEED_TRANSLATION_TABLE[prevIndex][0];
  const freqHigh = HZ_TO_SPEED_TRANSLATION_TABLE[index][0];

  if (freqHigh === freqLow) {
    // Avoid division by zero if the frequencies are the same
    return speedLow;
  }

  const proportion = (freq - freqLow) / (freqHigh - freqLow);
  return 1.1*(speedLow + proportion * (speedHigh - speedLow));
}

const HZ_TO_SPEED_TRANSLATION_TABLE= [
  [0.0, 0.0],
  [0.2, 1.0],
  [0.4, 3.0],
  [0.6, 7.0],
  [0.8, 15.0],
  [1.0, 25.0],
  [1.2, 33.0],
  [1.4, 36.0],
  [1.6, 40.0],
  [1.8, 40.0],
  [6.0, 40.0],
  [2.2, 35.0],
  [2.4, 33.0],
  [2.6, 27.0],
  [2.8, 20.0],
  [3.0, 10.0],
  [3.0, 5.0],
  [Number.MAX_VALUE, 0.0],
]

