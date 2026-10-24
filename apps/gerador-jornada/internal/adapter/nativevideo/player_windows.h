#ifndef JORNADA_NATIVE_VIDEO_H
#define JORNADA_NATIVE_VIDEO_H
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif
// JornadaVideoSnapshot transporta somente estado técnico, sem chamar Go a partir de COM.
typedef struct JornadaVideoSnapshot { int started, ended, audio; int32_t error; } JornadaVideoSnapshot;
void *jornada_video_open(uintptr_t parent, const uint16_t *path, int32_t *error);
void jornada_video_snapshot(void *player, JornadaVideoSnapshot *snapshot);
int jornada_video_resize(void *player);
// jornada_video_fit calcula um encaixe proporcional inteiro que nunca ultrapassa o cliente nativo.
void jornada_video_fit(int available_width, int available_height, int source_width, int source_height, int *rect);
// jornada_video_display_size converte razão de exibição (ex.: 16/9) em dimensões de quadro, não em pixels literais.
void jornada_video_display_size(int width, int height, int aspect_width, int aspect_height, int *size);
void jornada_video_close(void *player);
int jornada_window_alpha(uintptr_t window, int alpha);
#ifdef __cplusplus
}
#endif
#endif
