#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

uniform float x;
uniform float y;

void main() {
	float r = TexCoord.x * x;
	float b = TexCoord.y * y;

	float g = r*b-tan(r + b + 0.25*r*b);

	FragColor = vec4(r,g,b,1.0f);
}
