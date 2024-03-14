package com.fz.eureka.eurekafz;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.netflix.eureka.server.EnableEurekaServer;

@EnableEurekaServer
@SpringBootApplication
public class EurekaFzApplication {

	public static void main(String[] args) {
		SpringApplication.run(EurekaFzApplication.class, args);
	}

}
