/**
 * @fileoverview SSE (Server Sent Events) connection factory.
 * Uses the browser's native EventSource implementation.
 */

function getChannelName(url) {
  return url.split('/').filter(Boolean).pop() || url;
}

export function createSSEConnection(url, onMessage, onError, onOpen) {
  const channel = getChannelName(url);
  let manuallyClosed = false;
  let hasLoggedMessage = false;

  console.log(`[SSE] Opening connection: ${url}`);

  const es = new EventSource(url);

  es.onopen = () => {
    console.log(`[SSE] Connected: ${url}`);
    if (onOpen) onOpen();
  };

  es.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (!hasLoggedMessage) {
        console.log(`[SSE] Message received: ${channel}`);
        hasLoggedMessage = true;
      }
      onMessage(data);
    } catch (parseError) {
      console.error(`[SSE] Error`, parseError);
    }
  };

  es.onerror = (event) => {
    if (manuallyClosed) return;

    if (es.readyState === EventSource.CLOSED) {
      console.error(`[SSE] Error`, event);
      if (onError) onError(new Error(`EventSource closed for ${channel}`));
    }
  };

  return () => {
    if (manuallyClosed) return;
    manuallyClosed = true;
    es.close();
    console.log('[SSE] Connection closed manually');
  };
}
