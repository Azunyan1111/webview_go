#include "webview.h"

#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>

struct binding_context {
    webview_t w;
    uintptr_t index;
};

struct cookie_context {
    uintptr_t index;
};

void _webviewDispatchGoCallback(void *);
void _webviewBindingGoCallback(webview_t, char *, char *, uintptr_t);
void _webviewCookieGoCallback(char *, uintptr_t);
void _webviewClearCookiesGoCallback(int, uintptr_t);
void _webviewSetCookieGoCallback(int, uintptr_t);

static void _webview_dispatch_cb(webview_t w, void *arg) {
    _webviewDispatchGoCallback(arg);
}

static void _webview_binding_cb(const char *id, const char *req, void *arg) {
    struct binding_context *ctx = (struct binding_context *) arg;
    _webviewBindingGoCallback(ctx->w, (char *)id, (char *)req, ctx->index);
}

void CgoWebViewDispatch(webview_t w, uintptr_t arg) {
    webview_dispatch(w, _webview_dispatch_cb, (void *)arg);
}

void CgoWebViewBind(webview_t w, const char *name, uintptr_t index) {
    struct binding_context *ctx = calloc(1, sizeof(struct binding_context));
    ctx->w = w;
    ctx->index = index;
    webview_bind(w, name, _webview_binding_cb, (void *)ctx);
}

void CgoWebViewUnbind(webview_t w, const char *name) {
    webview_unbind(w, name);
}

static void _webview_cookie_cb(const char *cookies, void *arg) {
    struct cookie_context *ctx = (struct cookie_context *) arg;
    _webviewCookieGoCallback((char *)cookies, ctx->index);
    free(ctx);
}

void CgoWebViewGetCookies(webview_t w, uintptr_t index) {
    struct cookie_context *ctx = calloc(1, sizeof(struct cookie_context));
    ctx->index = index;
    webview_get_cookies(w, _webview_cookie_cb, (void *)ctx);
}

static void _webview_clear_cookies_cb(int success, void *arg) {
    struct cookie_context *ctx = (struct cookie_context *) arg;
    _webviewClearCookiesGoCallback(success, ctx->index);
    free(ctx);
}

void CgoWebViewClearCookies(webview_t w, uintptr_t index) {
    struct cookie_context *ctx = calloc(1, sizeof(struct cookie_context));
    ctx->index = index;
    webview_clear_cookies(w, _webview_clear_cookies_cb, (void *)ctx);
}

static void _webview_set_cookie_cb(int success, void *arg) {
    struct cookie_context *ctx = (struct cookie_context *) arg;
    _webviewSetCookieGoCallback(success, ctx->index);
    free(ctx);
}

void CgoWebViewSetCookie(webview_t w, const char *cookieJSON, uintptr_t index) {
    printf("CgoWebViewSetCookie: entered with index=%lu\n", (unsigned long)index);
    printf("CgoWebViewSetCookie: cookieJSON=%s\n", cookieJSON);
    struct cookie_context *ctx = calloc(1, sizeof(struct cookie_context));
    ctx->index = index;
    printf("CgoWebViewSetCookie: calling webview_set_cookie\n");
    webview_set_cookie(w, cookieJSON, _webview_set_cookie_cb, (void *)ctx);
    printf("CgoWebViewSetCookie: returned from webview_set_cookie\n");
}
