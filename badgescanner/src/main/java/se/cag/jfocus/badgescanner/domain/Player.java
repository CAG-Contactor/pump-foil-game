package se.cag.jfocus.badgescanner.domain;

import lombok.Builder;
import lombok.extern.slf4j.Slf4j;

import java.util.Objects;

@Builder(toBuilder = true)
@Slf4j
public record Player(String email, String name, String company, String fullString) {
    public boolean isValid() {
        boolean valid = Objects.nonNull(email) && ! email.isBlank() &&
                Objects.nonNull(name) && ! name.isBlank();
        if (! valid) {
            log.warn("User not valid: " + fullString());
        }
        return valid;
    }
}
