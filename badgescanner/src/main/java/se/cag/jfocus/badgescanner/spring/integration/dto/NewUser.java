package se.cag.jfocus.badgescanner.spring.integration.dto;


import lombok.Builder;

@Builder(toBuilder=true)
public record NewUser(
        String userId,
        String displayName,
        String organisation,
        String password)
{
}
