import io.undertow.Undertow;
import io.undertow.util.Headers;
import org.apache.commons.lang3.SystemUtils;

import static io.undertow.UndertowLogger.ROOT_LOGGER;

public class Handler {

    public static final String DEFAULT_PORT = "8080";

    public static void main(String[] args) {
        String port = SystemUtils.getEnvironmentVariable("PORT", DEFAULT_PORT);
        ROOT_LOGGER.infof("HTTP server listening on 0.0.0.0:" + port);
        Undertow server = Undertow.builder()
                .addHttpListener(Integer.parseInt(port), "0.0.0.0")
                .setHandler(exchange -> {
                    exchange.getResponseHeaders().put(Headers.CONTENT_TYPE, "application/json");
                    exchange.getResponseSender().send("""
                            {
                                "message":  "Hello, World from Scaleway Container !"
                            }""");
                }).build();
        server.start();
    }
}
