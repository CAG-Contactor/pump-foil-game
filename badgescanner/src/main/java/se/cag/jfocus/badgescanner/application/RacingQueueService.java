package se.cag.jfocus.badgescanner.application;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import se.cag.jfocus.badgescanner.annotations.JfocusService;

@JfocusService
@RequiredArgsConstructor
@Slf4j
public class RacingQueueService {
    private final RacingQueueRepo racingQueueRepo;
    private final ScanPlayerService scanPlayerService;

    public void enqueueNewUserByQr() {
        scanPlayerService.scanPlayer()
                .ifPresent(racingQueueRepo::enqueuePlayer);
    }
}
