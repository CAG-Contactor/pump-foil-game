package se.cag.jfocus.badgescanner.spring.integration.dto;


import lombok.Builder;

@Builder(toBuilder=true)
public record BadgeDTO(
        String email,
        String name,
        String company)
{
}
