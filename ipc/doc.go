/*
Package ipc provides a Client for connecting to the Sway window manager ipc
socket. The Client supports sending messages and subscribing to events.

Subsciption wraps Client to add support for typed event callbacks.

swager/ipc aims to be a fully featured library that supports all features
exposed over the sway ipc socket.

Notable missing pieces include everything related to Bars and Inputs, though the
primitives provided by the library should be able to get raw json representations.

Example usage can be seen in the swager/internal/* packages and the swager/blocks package.
*/
package ipc
