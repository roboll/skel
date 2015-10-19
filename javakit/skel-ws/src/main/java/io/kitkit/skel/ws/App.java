package io.kitkit.{{ .name }}.ws;

import io.dropwizard.Application;
import io.dropwizard.Configuration;
import io.dropwizard.setup.Bootstrap;
import io.dropwizard.setup.Environment;
import io.kitkit.wskit.WSKitBundle;


public class App extends Application<Configuration> {

    public static void main(String[] args) throws Exception {
        new App().run(args);
    }


    @Override
    public void initialize(Bootstrap<Configuration> bootstrap) {
        bootstrap.addBundle(new WSKitBundle());
    }

    @Override
    public void run(Configuration configuration, Environment environment) throws Exception {}
}
