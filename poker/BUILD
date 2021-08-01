load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "poker",
    srcs = [
        "game.go", 
        "cardset.go",
    ],
    importpath = "gopoker/poker",
    deps = [],
    visibility = ["//visibility:public"],
)

go_test(
    name = "cardset_test",
    srcs = [
        "cardset_test.go"
    ],
    embed = [
        ":poker",
    ],
)