from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.grpc import GrpcInstrumentorServer
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.instrumentation.grpc import server_interceptor

class OpenTelemetryManager:
    def __init__(self):
        self.trace_provider = TracerProvider()
        trace.set_tracer_provider(self.trace_provider)
        self.tracer = trace.get_tracer(__name__)

        otlp_exporter = OTLPSpanExporter()
        self.trace_provider.add_span_processor(BatchSpanProcessor(otlp_exporter))

        self.meter_provider = MeterProvider()
        metric_exporter = OTLPMetricExporter()
        metric_reader = PeriodicExportingMetricReader(metric_exporter)
        self.meter_provider._all_metric_readers = metric_reader
        self.meter = self.meter_provider.get_meter(__name__)

        self.embedding_time_metric = self.meter.create_histogram("embedding_creation_time", description="Time taken to create embeddings")

        GrpcInstrumentorServer().instrument()

    def start_trace(self, span_name):
        return self.tracer.start_as_current_span(span_name)

    def record_metric(self, start_time, end_time):
        self.embedding_time_metric.record(end_time - start_time)