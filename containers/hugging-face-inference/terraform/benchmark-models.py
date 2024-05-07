import csv
import json

import matplotlib.pyplot as plt
import pandas
import requests


class Benchmark:
    _model_families = ["llama", "mistral", "phi"]
    _endpoints = {}

    def __init__(
        self, models_file: str, benchmark_file: str, results_figure: str, message: str
    ) -> None:
        self.models_file = models_file
        self.benchmark_file = benchmark_file
        self.message = message
        self.results_figure = results_figure

    def get_container_endpoints_from_json_file(self) -> None:
        if self.models_file == "":
            raise Exception("file name is empty")

        with open(self.models_file, "r") as models_file:
            json_data = json.load(models_file)

        for family in self._model_families:
            self._endpoints[family] = []
            for model in json_data[family]:
                self._endpoints[family].append(
                    {"model": model["file"], "endpoint": model["ctn_endpoint"]}
                )

    def analyze_results(self) -> None:
        benchmark_results = pandas.read_csv(self.benchmark_file)
        benchmark_results.boxplot(column="Total Response Time", by="Family").plot()
        plt.ylabel("Total Response Time in seconds")
        plt.savefig(self.results_figure)

    def benchmark_models(self, num_samples: int) -> None:
        self.get_container_endpoints_from_json_file()

        fields = ["Model", "Family", "Total Response Time", "Response Message"]
        benchmark_data = []

        for family in self._model_families:
            for endpoint in self._endpoints[family]:
                if endpoint["endpoint"] == "":
                    raise Exception("model endpoint is empty")

                for _ in range(num_samples):
                    try:
                        print(
                            "Calling model {model} on endpoint {endpoint} with message {message}".format(
                                model=endpoint["model"],
                                endpoint=endpoint["endpoint"],
                                message=self.message,
                            )
                        )

                        rsp = requests.post(
                            endpoint["endpoint"], json={"message": self.message}
                        )

                        response_text = rsp.json()["choices"][0]["text"]

                        print(
                            "The model {model} responded with: {response_text}".format(
                                model=endpoint["model"], response_text=response_text
                            )
                        )

                        benchmark_data.append(
                            [
                                endpoint["model"],
                                family,
                                rsp.elapsed.total_seconds(),
                                response_text,
                            ]
                        )
                    except:
                        pass

        with open(self.benchmark_file, "w") as results_file:
            wrt = csv.writer(results_file)
            wrt.writerow(fields)
            wrt.writerows(benchmark_data)

        self.analyze_results()


if __name__ == "__main__":

    benchmark = Benchmark(
        models_file="hf-models.json",
        benchmark_file="benchmark-results.csv",
        results_figure="results-plot.png",
        message="What the difference between an elephant and an ant?",
    )

    benchmark.benchmark_models(num_samples=50)
