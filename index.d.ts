import * as z from "zod";

declare global {
  interface HTMXBeforeSwapEvent extends Event {
    detail: {
      xhr: {
        status: number;
      };
      shouldSwap: boolean;
      isError: boolean;
    };
  }

  interface HTMLElementEventMap {
    ["htmx:beforeSwap"]: HTMXBeforeSwapEvent;
  }

  type HttpMethod = "GET" | "POST" | "PATCH" | "PUT" | "DELETE";

  interface HTMX {
    ajax(
      method: HttpMethod,
      url: string,
      selector: HTMLElement
    ): Promise<unknown>;

    ajax(
      method: HttpMethod,
      url: string,
      selector: string
    ): Promise<unknown>;

    ajax(
      method: HttpMethod,
      url: string,
      context: HtmxAjaxHelperContext
    ): Promise<unknown>;
  }
  interface HtmxAjaxHelperContext {
    source?: Element | string;
    event?: Event;
    handler?: HtmxAjaxHandler;
    target?: Element | string;
    swap?: HtmxSwapStyle;
    values?: Object | FormData;
    headers?: Record<string, string>;
    select?: string;
  }

  type HtmxSwapStyle = "innerHTML" | "outerHTML" | "beforebegin" | "afterbegin" | "beforeend" | "afterend" | "delete" | "none" | string;

  type HtmxAjaxHandler = (element: Element) => HtmxResponseInfo;

  interface HtmxResponseInfo {
    xhr: XMLHttpRequest;
    target: Element;
    requestConfig: unknown;
    etc: unknown;
    boosted: boolean;
    select: string;
    pathInfo: { requestPath: string, finalRequestPath: string, responsePath: string | null, anchor: string };
    failed?: boolean;
    successful?: boolean;
    keepIndicators?: boolean;
  }
  interface Window {
    Zod: typeof z;
    htmx: HTMX;
  }
}

export { };
