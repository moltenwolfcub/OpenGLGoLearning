#version 330 core
out vec4 FragColor;

in vec3 ModelPos;

void main() {
	float scalar = abs(ModelPos.x/5);

	FragColor = vec4(0.28*scalar, 0.38*scalar, 0.94*scalar, 1.0);
}
