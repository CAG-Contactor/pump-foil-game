package se.cag.jfocus.badgescanner.application;


import se.cag.jfocus.badgescanner.domain.Player;

public interface RacingQueueRepo {

    void enqueuePlayer(Player player);

}
