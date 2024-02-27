import asyncio
import builtins
from asyncio import AbstractEventLoop
from contextlib import ContextDecorator
from typing import Any, Coroutine, Union

from grpclib.client import Channel

from blobcache import BlobCacheStub, StoreContentRequest


class InterceptReadFile:
    def __init__(self, file, blob_cache_stub):
        self._file = file
        self._blob_cache_stub: BlobCacheStub = blob_cache_stub

    def read(self, *args, **kwargs) -> Any:
        content = self._file.read(*args, **kwargs)

        if isinstance(content, str):
            content = content.encode("utf-8")

        r = Cache.run_sync(
            self._blob_cache_stub.store_content(
                store_content_request_iterator=[StoreContentRequest(content=content)]
            )
        )
        print(r)

        return content

    def __getattr__(self, name):
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
            file = self.original_open(*args, **kwargs)  # Use the original open

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
