import asyncio
import builtins
from asyncio import AbstractEventLoop
from contextlib import ContextDecorator
from typing import Any, Coroutine, Union

from grpclib.client import Channel

from blobcache import BlobCacheStub, StoreContentRequest


class InterceptReadFile:
    def __init__(self, file, blob_cache_stub: BlobCacheStub):
        self._file = file
        self._blob_cache_stub = blob_cache_stub

    def read(self, size=-1) -> Any:
        # If a specific size is not requested, use generator for chunked read
        if size <= 0:
            return self._read_and_store_in_chunks()
        else:
            chunk = self._file.read(size)
            # Ideally, you should also handle this single read with BlobCache
            return chunk

    def _read_and_store_in_chunks(self, chunk_size=1024 * 1024) -> bytes:
        """A generator that yields StoreContentRequest for each chunk read."""

        def generate_chunks():
            i = 0
            while True:
                chunk = self._file.read(chunk_size)
                print(i, ":", chunk)

                if not chunk:
                    break  # End of file

                yield StoreContentRequest(content=chunk)

                i += 1

        # Asynchronously send chunks to BlobCache
        r = Cache.run_sync(
            self._blob_cache_stub.store_content(
                store_content_request_iterator=generate_chunks()
            )
        )
        print(r)

    def __getattr__(self, name):
        # Delegate attribute access to the original file object
        return getattr(self._file, name)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self._file.close()


class Cache(ContextDecorator):
    def __init__(self, host: str = "localhost", port: int = 2049, ssl: bool = False):
        self.channel = Channel(host=host, port=port, ssl=ssl)
        self.blob_cache_stub: BlobCacheStub = BlobCacheStub(channel=self.channel)
        self.original_open = open

    def __enter__(self):
        def custom_open(*args, **kwargs):
            file = self.original_open(*args, "rb")  # Use the original open

            return InterceptReadFile(
                file, self.blob_cache_stub
            )  # Return a custom wrapped file object

        builtins.open = custom_open  # Patch the built-in open
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        builtins.open = self.original_open  # Restore the original open
        self.channel.close()

    @staticmethod
    def run_sync(
        coroutine: Coroutine, loop: Union[AbstractEventLoop, None] = None
    ) -> Any:
        if loop is None:
            loop = asyncio.get_event_loop()

        return loop.run_until_complete(coroutine)
